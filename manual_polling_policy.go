package configcat

// Describes a RefreshPolicy which fetches the latest configuration over HTTP every time when a get configuration is called.
type ManualPollingPolicy struct {
	ConfigRefresher
}

// Initializes a new ManualPollingPolicy.
func NewManualPollingPolicy(
	fetcher ConfigProvider,
	store *ConfigStore) *ManualPollingPolicy {
	return &ManualPollingPolicy{ ConfigRefresher: ConfigRefresher{ Fetcher:fetcher, Store:store }}
}

// GetConfigurationAsync reads the current configuration value.
func (policy *ManualPollingPolicy) GetConfigurationAsync() *AsyncResult {
	return policy.Fetcher.GetConfigurationAsync().ApplyThen(func(result interface{}) interface{} {
		response := result.(FetchResponse)

		cached := policy.Store.Get()
		if response.IsFetched() {
			fetched := response.Body
			if cached != fetched {
				policy.Store.Set(fetched)
			}

			return fetched
		}

		return cached
	})
}

// Close shuts down the policy.
func (policy *ManualPollingPolicy) Close() {
}
