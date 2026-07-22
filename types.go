package capydbclient

import (
	"encoding/json"
	"time"
)

// This file holds the canonical Go representations of the control-plane entities
// and request bodies that the CLI and the Terraform provider both consume. They
// mirror the component schemas in backend/internal/httpapi/openapi.json (the
// same source the TypeScript SDK is generated from) so the two Go clients share
// one shape instead of each re-declaring — and drifting — their own copies.
//
// Response entities carry the full schema field set; a consumer that only reads
// a subset simply ignores the rest. Date-time fields use time.Time (required)
// or *time.Time (optional), matching the backend model.

// Organization is a billing/identity tenant.
type Organization struct {
	BillingCustomerID     string     `json:"billing_customer_id,omitempty"`
	BillingEmail          string     `json:"billing_email,omitempty"`
	BillingName           string     `json:"billing_name,omitempty"`
	BillingPeriodEnd      *time.Time `json:"billing_period_end,omitempty"`
	BillingPlan           string     `json:"billing_plan"`
	BillingProductID      string     `json:"billing_product_id,omitempty"`
	BillingProvider       string     `json:"billing_provider"`
	BillingStatus         string     `json:"billing_status"`
	BillingSubscriptionID string     `json:"billing_subscription_id,omitempty"`
	ClerkOrganizationID   string     `json:"clerk_organization_id,omitempty"`
	ClerkOrganizationSlug string     `json:"clerk_organization_slug,omitempty"`
	CreatedAt             time.Time  `json:"created_at"`
	ID                    string     `json:"id"`
	Name                  string     `json:"name"`
	Slug                  string     `json:"slug"`
	SuspendedAt           *time.Time `json:"suspended_at,omitempty"`
	SuspendedReason       string     `json:"suspended_reason,omitempty"`
	UpdatedAt             time.Time  `json:"updated_at"`
	VercelInstallationID  string     `json:"vercel_installation_id,omitempty"`
}

// Viewer is the GET /v1/me organization payload.
type Viewer struct {
	Organization *Organization `json:"organization"`
}

// Project is a managed Postgres database.
type Project struct {
	CreatedAt              time.Time `json:"created_at"`
	DatabaseName           string    `json:"database_name"`
	DirectPort             int       `json:"direct_port"`
	Environment            string    `json:"environment"`
	ID                     string    `json:"id"`
	IdleTransactionTimeout string    `json:"idle_transaction_timeout"`
	LastError              string    `json:"last_error,omitempty"`
	LatestJobID            string    `json:"latest_job_id,omitempty"`
	MaxConnections         int       `json:"max_connections"`
	Name                   string    `json:"name"`
	OrganizationID         string    `json:"organization_id"`
	Plan                   string    `json:"plan"`
	PooledPort             int       `json:"pooled_port"`
	// PostgresVersion is the major version of the project's database. Empty
	// while the database is still provisioning.
	PostgresVersion   string `json:"postgres_version,omitempty"`
	PrimaryInstanceID string `json:"primary_instance_id,omitempty"`
	PublicHost        string `json:"public_host,omitempty"`
	Region            string `json:"region"`
	RoleName          string `json:"role_name"`
	// RuntimeStatus overlays the scale-to-zero lifecycle in customer
	// vocabulary: "active", "paused", or "resuming". Empty while the database
	// is provisioning or being deleted (see State).
	RuntimeStatus     string    `json:"runtime_status,omitempty"`
	Slug              string    `json:"slug"`
	SSLMode           string    `json:"ssl_mode,omitempty"`
	State             string    `json:"state"`
	StatementTimeout  string    `json:"statement_timeout"`
	StorageLimitBytes int64     `json:"storage_limit_bytes"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// Job is an asynchronous control-plane operation. Poll until State is
// "completed" or "failed".
type Job struct {
	Attempts            int        `json:"attempts"`
	ClaimedAt           *time.Time `json:"claimed_at,omitempty"`
	ClaimedBy           string     `json:"claimed_by,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	Error               string     `json:"error,omitempty"`
	HostID              string     `json:"host_id,omitempty"`
	ID                  string     `json:"id"`
	InstanceID          string     `json:"instance_id,omitempty"`
	LastExitCode        int        `json:"last_exit_code,omitempty"`
	LastStderr          string     `json:"last_stderr,omitempty"`
	LastStdout          string     `json:"last_stdout,omitempty"`
	LockedResource      string     `json:"locked_resource,omitempty"`
	MaxAttempts         int        `json:"max_attempts"`
	OrganizationID      string     `json:"organization_id"`
	PreviewDatabaseID   string     `json:"preview_database_id,omitempty"`
	ProjectID           string     `json:"project_id,omitempty"`
	RetryClassification string     `json:"retry_classification,omitempty"`
	// Result is the structured result payload of a completed job. The API
	// returns it only for job types whose result is part of the product
	// contract (currently project.import_follow_status).
	Result    json.RawMessage `json:"result,omitempty"`
	StartedAt *time.Time      `json:"started_at,omitempty"`
	State     string          `json:"state"`
	Type      string          `json:"type"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ConnectionInfo is a project or preview database's connection endpoints.
type ConnectionInfo struct {
	DirectURL string `json:"direct_url,omitempty"`
	PooledURL string `json:"pooled_url,omitempty"`
	Username  string `json:"username"`
}

// APIKey is an organization or project-scoped API key. Plaintext secrets are
// returned only on creation, never on list endpoints.
type APIKey struct {
	CreatedAt       time.Time  `json:"created_at"`
	CreatedByUserID string     `json:"created_by_user_id,omitempty"`
	DeviceName      string     `json:"device_name,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	ID              string     `json:"id"`
	IsActive        bool       `json:"is_active"`
	KeyPrefix       string     `json:"key_prefix"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
	Name            string     `json:"name"`
	OrganizationID  string     `json:"organization_id"`
	ProjectID       string     `json:"project_id,omitempty"`
	RevokedAt       *time.Time `json:"revoked_at,omitempty"`
	Scopes          []string   `json:"scopes"`
	Source          string     `json:"source"`
}

// WebhookEndpoint is an outbound webhook receiver.
type WebhookEndpoint struct {
	CreatedAt      time.Time `json:"created_at"`
	Description    string    `json:"description,omitempty"`
	EventTypes     []string  `json:"event_types"`
	ID             string    `json:"id"`
	IsActive       bool      `json:"is_active"`
	OrganizationID string    `json:"organization_id"`
	UpdatedAt      time.Time `json:"updated_at"`
	URL            string    `json:"url"`
}

// CreateProjectRequest is the POST /v1/projects body.
type CreateProjectRequest struct {
	Environment    string `json:"environment,omitempty"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id,omitempty"`
	// PostgresVersion picks the database's major version ("16", "17", "18").
	// Omit for the platform default. Immutable after creation; previews and
	// restores inherit it.
	PostgresVersion string `json:"postgres_version,omitempty"`
	Region          string `json:"region,omitempty"`
	Slug            string `json:"slug,omitempty"`
}

// CreatePreviewRequest is the create-preview-database body.
type CreatePreviewRequest struct {
	Mode     string `json:"mode,omitempty"`
	Name     string `json:"name,omitempty"`
	TTLHours int    `json:"ttl_hours,omitempty"`
}

// CreateAPIKeyRequest is the create-API-key body.
type CreateAPIKeyRequest struct {
	DeviceName string     `json:"device_name,omitempty"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Name       string     `json:"name"`
	ProjectID  string     `json:"project_id,omitempty"`
	Scopes     []string   `json:"scopes"`
}

// CreateWebhookEndpointRequest is the create-webhook-endpoint body.
type CreateWebhookEndpointRequest struct {
	Description string   `json:"description,omitempty"`
	EventTypes  []string `json:"event_types,omitempty"`
	URL         string   `json:"url"`
}

// DatabaseSchema is the canonical schema document from GET
// /v1/projects/{id}/schema (and its preview-database sibling): one
// introspection pass covering schemas, tables, views, columns, constraints,
// enums and extensions. It is the input contract for type generation and
// schema diffing.
type DatabaseSchema struct {
	DatabaseName    string            `json:"database_name"`
	Extensions      []SchemaExtension `json:"extensions"`
	PostgresVersion string            `json:"postgres_version"`
	Schemas         []SchemaNamespace `json:"schemas"`
}

// SchemaExtension is an installed Postgres extension.
type SchemaExtension struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// SchemaNamespace is one Postgres schema (namespace) with its relations and
// enum types.
type SchemaNamespace struct {
	Enums  []SchemaEnum  `json:"enums"`
	Name   string        `json:"name"`
	Tables []SchemaTable `json:"tables"`
}

// SchemaEnum is a user-defined enum type.
type SchemaEnum struct {
	Comment string   `json:"comment,omitempty"`
	Name    string   `json:"name"`
	Values  []string `json:"values"`
}

// SchemaTable is a relation; Kind is one of table, partitioned_table, view,
// materialized_view, foreign_table.
type SchemaTable struct {
	Columns           []SchemaColumn           `json:"columns"`
	Comment           string                   `json:"comment,omitempty"`
	ForeignKeys       []SchemaForeignKey       `json:"foreign_keys"`
	Kind              string                   `json:"kind"`
	Name              string                   `json:"name"`
	PrimaryKey        []string                 `json:"primary_key"`
	UniqueConstraints []SchemaUniqueConstraint `json:"unique_constraints"`
}

// SchemaColumn is one column of a relation. For array columns UDTName and
// UDTSchema describe the element type.
type SchemaColumn struct {
	// ArrayDims is the declared array dimension count (>= 1 when IsArray).
	ArrayDims   int    `json:"array_dims,omitempty"`
	Comment     string `json:"comment,omitempty"`
	DataType    string `json:"data_type"`
	Default     string `json:"default,omitempty"`
	Identity    string `json:"identity,omitempty"`
	IsArray     bool   `json:"is_array"`
	IsEnum      bool   `json:"is_enum"`
	IsGenerated bool   `json:"is_generated"`
	IsNullable  bool   `json:"is_nullable"`
	Name        string `json:"name"`
	Position    int    `json:"position"`
	UDTName     string `json:"udt_name"`
	UDTSchema   string `json:"udt_schema"`
}

// SchemaForeignKey is a foreign-key constraint.
type SchemaForeignKey struct {
	Columns           []string `json:"columns"`
	Name              string   `json:"name"`
	OnDelete          string   `json:"on_delete"`
	OnUpdate          string   `json:"on_update"`
	ReferencedColumns []string `json:"referenced_columns"`
	ReferencedSchema  string   `json:"referenced_schema"`
	ReferencedTable   string   `json:"referenced_table"`
}

// SchemaUniqueConstraint is a unique constraint.
type SchemaUniqueConstraint struct {
	Columns []string `json:"columns"`
	Name    string   `json:"name"`
}

// GeneratedTypes is one generated source file from GET
// /v1/projects/{id}/schema/types: TypeScript interfaces, Zod schemas or a
// Drizzle schema rendered server-side from the live database schema.
type GeneratedTypes struct {
	Content  string `json:"content"`
	Filename string `json:"filename"`
	Language string `json:"language"`
	Style    string `json:"style,omitempty"`
}
