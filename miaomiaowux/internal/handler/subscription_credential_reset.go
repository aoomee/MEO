package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"

	"miaomiaowux/internal/storage"
)

type subscriptionCredentialResetResult struct {
	Token              string
	CredentialsUpdated int
	NodesUpdated       int
}

type subscriptionCredentialRecord struct {
	kind         string
	id           int64
	serverID     int64
	inboundTag   string
	protocol     string
	email        string
	old          map[string]interface{}
	oldJSON      string
	newJSON      string
	active       bool
	routedNodeID int64
}

// resetUserSubscriptionCredentials rotates both sides of a leaked
// subscription: every Xray client credential associated with the current user
// and every URL credential that can produce a subscription. Remote Xray
// changes use the existing prepare/commit transaction; database references and
// short codes are committed atomically afterwards.
func resetUserSubscriptionCredentials(ctx context.Context, repo *storage.TrafficRepository, rm *RemoteManageHandler, user storage.User) (subscriptionCredentialResetResult, error) {
	var result subscriptionCredentialResetResult
	if rm == nil {
		return result, errors.New("remote management is unavailable")
	}
	records, err := collectSubscriptionCredentialRecords(ctx, repo, user.Username)
	if err != nil {
		return result, err
	}

	newByRemoteKey := make(map[string]map[string]interface{})
	for i := range records {
		key := xrayCredentialRemoteKey(records[i].serverID, records[i].inboundTag, records[i].protocol, records[i].old)
		rotated := newByRemoteKey[key]
		if rotated == nil {
			rotated, err = rotateXrayCredential(records[i].protocol, records[i].old)
			if err != nil {
				return result, fmt.Errorf("rotate %s credential %d: %w", records[i].kind, records[i].id, err)
			}
			newByRemoteKey[key] = rotated
		}
		encoded, marshalErr := json.Marshal(rotated)
		if marshalErr != nil {
			return result, marshalErr
		}
		records[i].newJSON = string(encoded)
	}

	states := map[int64]*packageConfigState{}
	changedServers := map[int64]bool{}
	processedRemote := map[string]bool{}
	for _, record := range records {
		if !record.active {
			continue
		}
		key := xrayCredentialRemoteKey(record.serverID, record.inboundTag, record.protocol, record.old)
		if processedRemote[key] {
			continue
		}
		processedRemote[key] = true
		state, loadErr := loadPackageServerConfig(ctx, rm, states, record.serverID)
		if loadErr != nil {
			return result, fmt.Errorf("load server %d Xray config: %w", record.serverID, loadErr)
		}
		replaced, replaceErr := replaceManagedClientCredential(state.config, record.inboundTag, record.old, newByRemoteKey[key])
		if replaceErr != nil {
			return result, fmt.Errorf("replace %s credential %d: %w", record.kind, record.id, replaceErr)
		}
		if replaced {
			changedServers[record.serverID] = true
		}
	}

	changes, err := buildSubscriptionCredentialDBChanges(ctx, repo, user, records)
	if err != nil {
		return result, err
	}
	targets := make([]xrayConfigTransactionTarget, 0, len(changedServers))
	for serverID := range changedServers {
		state := states[serverID]
		if err := validateManagedConfig(state.config); err != nil {
			return result, fmt.Errorf("validate server %d Xray config: %w", serverID, err)
		}
		configJSON, err := json.MarshalIndent(state.config, "", "  ")
		if err != nil {
			return result, err
		}
		targets = append(targets, xrayConfigTransactionTarget{ServerID: serverID, Config: string(configJSON)})
	}

	operationID := "credential-reset-" + uuid.NewString()
	if len(targets) > 0 {
		if err := applyXrayConfigTransaction(ctx, rm, operationID, targets); err != nil {
			return result, err
		}
	}
	token, err := repo.CommitSubscriptionCredentialRotation(ctx, user.Username, changes)
	if err != nil {
		if len(targets) > 0 {
			_ = finishXrayConfigTransaction(context.WithoutCancel(ctx), rm, operationID, targets, false)
		}
		return result, fmt.Errorf("commit subscription credential rotation: %w", err)
	}
	if len(targets) > 0 {
		if err := finishXrayConfigTransaction(context.WithoutCancel(ctx), rm, operationID, targets, true); err != nil {
			log.Printf("[CredentialReset] operation %s committed but agent cleanup incomplete: %v", operationID, err)
		}
		for _, target := range targets {
			go rm.refreshXraySnapshot(target.ServerID)
		}
	}
	result.Token = token
	result.CredentialsUpdated = len(newByRemoteKey)
	result.NodesUpdated = len(changes.Nodes)
	return result, nil
}

func collectSubscriptionCredentialRecords(ctx context.Context, repo *storage.TrafficRepository, username string) ([]subscriptionCredentialRecord, error) {
	var records []subscriptionCredentialRecord
	appendRecord := func(record subscriptionCredentialRecord, raw string) error {
		if err := json.Unmarshal([]byte(raw), &record.old); err != nil || record.old == nil {
			return fmt.Errorf("decode %s credential %d: %w", record.kind, record.id, err)
		}
		record.oldJSON = raw
		records = append(records, record)
		return nil
	}
	legacyInbound, err := repo.GetUserInboundConfigs(ctx, username)
	if err != nil {
		return nil, err
	}
	for _, cfg := range legacyInbound {
		if err := appendRecord(subscriptionCredentialRecord{kind: "legacy-inbound", id: cfg.ID, serverID: cfg.ServerID, inboundTag: cfg.InboundTag, protocol: cfg.Protocol, active: true}, cfg.CredentialJSON); err != nil {
			return nil, err
		}
	}
	assignmentInbound, err := repo.ListPackageAssignmentInboundConfigsByUser(ctx, username)
	if err != nil {
		return nil, err
	}
	for _, cfg := range assignmentInbound {
		if err := appendRecord(subscriptionCredentialRecord{kind: "assignment-inbound", id: cfg.ID, serverID: cfg.ServerID, inboundTag: cfg.InboundTag, protocol: cfg.Protocol, email: cfg.Email, active: true}, cfg.CredentialJSON); err != nil {
			return nil, err
		}
	}

	serverByName := map[string]int64{}
	servers, err := repo.ListRemoteServers(ctx)
	if err != nil {
		return nil, err
	}
	for _, server := range servers {
		serverByName[server.Name] = server.ID
	}
	appendRouted := func(kind string, id, nodeID int64, email, raw string, active bool) error {
		node, err := repo.GetRoutedNodeDetail(ctx, nodeID)
		if err != nil {
			return fmt.Errorf("load routed node %d: %w", nodeID, err)
		}
		protocol := node.Protocol
		inboundTag := node.InboundTag
		serverName := node.OriginalServer
		if node.ParentNodeID != nil {
			if parent, parentErr := repo.GetNodeByID(ctx, *node.ParentNodeID); parentErr == nil {
				protocol, inboundTag, serverName = parent.Protocol, parent.InboundTag, parent.OriginalServer
			}
		}
		serverID := serverByName[serverName]
		if serverID <= 0 {
			return fmt.Errorf("routed node %d server %q not found", nodeID, serverName)
		}
		return appendRecord(subscriptionCredentialRecord{kind: kind, id: id, serverID: serverID, inboundTag: inboundTag, protocol: protocol, email: email, active: active, routedNodeID: nodeID}, raw)
	}
	legacySubaccounts, err := repo.ListUserSubaccounts(ctx, username)
	if err != nil {
		return nil, err
	}
	for _, sa := range legacySubaccounts {
		if err := appendRouted("legacy-routed", sa.ID, sa.RoutedNodeID, sa.Email, sa.CredentialJSON, sa.IsActive); err != nil {
			return nil, err
		}
	}
	assignmentSubaccounts, err := repo.ListPackageAssignmentSubaccountsByUser(ctx, username)
	if err != nil {
		return nil, err
	}
	for _, sa := range assignmentSubaccounts {
		if err := appendRouted("assignment-routed", sa.ID, sa.RoutedNodeID, sa.Email, sa.CredentialJSON, sa.IsActive); err != nil {
			return nil, err
		}
	}
	return records, nil
}

func replaceManagedClientCredential(config map[string]interface{}, inboundTag string, oldCredential, newCredential map[string]interface{}) (bool, error) {
	_, settings, protocol, err := managedInbound(config, inboundTag)
	if err != nil {
		return false, err
	}
	key, err := managedClientArrayKey(protocol)
	if err != nil {
		return false, err
	}
	entries, _ := settings[key].([]interface{})
	for i, raw := range entries {
		entry, _ := raw.(map[string]interface{})
		if entry != nil && matchCredential(entry, oldCredential, protocol) {
			entries[i] = cloneMap(newCredential)
			settings[key] = entries
			return true, nil
		}
	}
	// A disabled/over-limit account is intentionally absent from Xray. Its
	// stored credential still rotates, but the reset must not reactivate it.
	return false, nil
}

func buildSubscriptionCredentialDBChanges(ctx context.Context, repo *storage.TrafficRepository, user storage.User, records []subscriptionCredentialRecord) (storage.SubscriptionCredentialRotation, error) {
	changes := storage.SubscriptionCredentialRotation{}
	nodes, err := repo.ListAllNodes(ctx)
	if err != nil {
		return changes, err
	}
	servers, err := repo.ListRemoteServers(ctx)
	if err != nil {
		return changes, err
	}
	serverNameByID := map[int64]string{}
	for _, server := range servers {
		serverNameByID[server.ID] = server.Name
	}
	nodeChanges := map[int64]storage.NodeCredentialRotation{}
	suspensionChanges := map[string]string{}
	for _, record := range records {
		row := storage.CredentialJSONRotation{ID: record.id, CredentialJSON: record.newJSON}
		switch record.kind {
		case "legacy-inbound":
			changes.LegacyInbound = append(changes.LegacyInbound, row)
		case "assignment-inbound":
			changes.AssignmentInbound = append(changes.AssignmentInbound, row)
		case "legacy-routed":
			changes.LegacySubaccount = append(changes.LegacySubaccount, row)
		case "assignment-routed":
			changes.AssignmentSubaccount = append(changes.AssignmentSubaccount, row)
		}
		suspensionChanges[record.oldJSON] = record.newJSON
		if user.Role != storage.RoleAdmin {
			continue
		}
		var newCredential map[string]interface{}
		if err := json.Unmarshal([]byte(record.newJSON), &newCredential); err != nil {
			return changes, err
		}
		if record.routedNodeID > 0 {
			routed, err := repo.GetRoutedNodeDetail(ctx, record.routedNodeID)
			if err != nil || routed.RoutedAdminEmail == "" || routed.RoutedAdminEmail != record.email {
				continue
			}
			clash := cloneClashWithCredential(routed.ClashConfig, record.protocol, newCredential, routed.NodeName)
			nodeChanges[routed.ID] = storage.NodeCredentialRotation{ID: routed.ID, CredentialJSON: record.newJSON, ClashConfig: clash, RoutedAdmin: true}
			continue
		}
		serverName := serverNameByID[record.serverID]
		for _, node := range nodes {
			if node.NodeType == "routed" || (node.Username != "" && node.Username != user.Username) || node.OriginalServer != serverName || node.InboundTag != record.inboundTag || !clashUsesCredential(node.ClashConfig, record.protocol, record.old) {
				continue
			}
			clash := cloneClashWithCredential(node.ClashConfig, record.protocol, newCredential, node.NodeName)
			nodeChanges[node.ID] = storage.NodeCredentialRotation{ID: node.ID, ClashConfig: clash}
		}
	}
	for oldJSON, newJSON := range suspensionChanges {
		changes.Suspensions = append(changes.Suspensions, storage.CredentialJSONReplacement{Old: oldJSON, New: newJSON})
	}
	for _, change := range nodeChanges {
		changes.Nodes = append(changes.Nodes, change)
	}
	return changes, nil
}
