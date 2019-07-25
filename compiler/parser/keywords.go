package parser

// Keyword represents a keyword
type Keyword = string

const (
	// KeywordSchema represents the 'schema' keyword
	KeywordSchema Keyword = "schema"

	// KeywordAlias represents the 'alias' keyword
	KeywordAlias Keyword = "alias"

	// KeywordUnion represents the 'union' keyword
	KeywordUnion Keyword = "union"

	// KeywordEnum represents the 'enum' keyword
	KeywordEnum Keyword = "enum"

	// KeywordResolver represents the 'resolver' keyword
	KeywordResolver Keyword = "resolver"

	// KeywordStruct represents the 'struct' keyword
	KeywordStruct Keyword = "struct"

	// KeywordTrait represents the 'trait' keyword
	KeywordTrait Keyword = "trait"

	// KeywordQuery represents the 'query' keyword
	KeywordQuery Keyword = "query"

	// KeywordMutation represents the 'mutation' keyword
	KeywordMutation Keyword = "mutation"

	// KeywordSubscription represents the 'subscription' keyword
	KeywordSubscription Keyword = "subscription"
)
