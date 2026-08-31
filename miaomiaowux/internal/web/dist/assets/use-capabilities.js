const features = [
  "limiter",
  "embedded",
  "speed_test",
  "custom_branding",
  "server_share",
  "reality_pool",
]

const plan = {
  name: "",
  display_name: "MEO",
  description: "MEO",
  max_servers: 999999,
  max_nodes: 999999,
  max_users: 999999,
  features,
}

const status = { valid: true, plan }

const queryResult = (data) => ({
  data,
  error: null,
  isLoading: false,
  isFetching: false,
  refetch: async () => ({ data }),
})

function useCapabilities() {
  return queryResult(status)
}

function useFeature() {
  return { hasFeature: true, plan }
}

function usesDecorativeTheme() {
  return false
}

function useCapabilityUsage() {
  return queryResult({
    success: true,
    server_count: 0,
    node_count: 0,
    user_count: 0,
    plan,
  })
}

export {
  useCapabilities as a,
  useFeature as b,
  usesDecorativeTheme as c,
  useCapabilityUsage as u,
}
