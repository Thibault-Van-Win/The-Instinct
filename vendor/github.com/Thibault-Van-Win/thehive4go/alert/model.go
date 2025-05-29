package alert

import (
	"github.com/Thibault-Van-Win/thehive4go/observable"
	"github.com/Thibault-Van-Win/thehive4go/procedure"
)

type Alert struct {
	ID              string           `json:"_id"`
	Type            string           `json:"type"`
	Source          string           `json:"source"`
	SourceRef       string           `json:"sourceRef"`
	ExternalLink    string           `json:"externalLink"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	Severity        int              `json:"severity"`
	SeverityLabel   string           `json:"severityLabel"`
	Date            int64            `json:"date"`
	Tags            []string         `json:"tags"`
	TLP             int              `json:"tlp"`
	TLPLabel        string           `json:"tlpLabel"`
	PAP             int              `json:"pap"`
	PAPLabel        string           `json:"papLabel"`
	Follow          bool             `json:"follow"`
	CustomFields    []map[string]any `json:"customFields"`
	CaseTemplate    string           `json:"caseTemplate"`
	ObservableCount int              `json:"observableCount"`
	CaseID          string           `json:"caseId"`
	Status          string           `json:"status"`
	Stage           string           `json:"stage"`
	Assignee        string           `json:"assignee"`
	Summary         string           `json:"summary"`
	ExtraData       map[string]any   `json:"extraData"`

	NewDate           int64 `json:"newDate"`
	InProgressDate    int64 `json:"inProgressDate"`
	ClosedDate        int64 `json:"closedDate"`
	ImportedDate      int64 `json:"importedDate"`
	TimeToDetect      int   `json:"timeToDetect"`
	TimeToTriage      int   `json:"timeToTriage"`
	TimeToQualify     int   `json:"timeToQualify"`
	TimeToAcknowledge int   `json:"timeToAcknowledge"`

	CreatedBy string `json:"_createdBy"`
	UpdatedBy string `json:"_updatedBy"`
	CreatedAt int64  `json:"_createdAt"`
	UpdatedAt int64  `json:"_updatedAt"`
}

type CreateAlertRequest struct {
	Type         string                         `json:"type"`
	Source       string                         `json:"source"`
	SourceRef    string                         `json:"sourceRef"`
	ExternalLink string                         `json:"externalLink,omitempty"`
	Title        string                         `json:"title"`
	Description  string                         `json:"description"`
	Severity     int                            `json:"severity,omitempty"`
	Date         int64                          `json:"date,omitempty"`
	Tags         []string                       `json:"tags,omitempty"`
	Flag         bool                           `json:"flag,omitempty"`
	TLP          int                            `json:"tlp,omitempty"`
	PAP          int                            `json:"pap,omitempty"`
	CustomFields map[string]any                 `json:"customFields,omitempty"`
	Summary      string                         `json:"summary,omitempty"`
	Status       string                         `json:"status,omitempty"`
	Assignee     string                         `json:"assignee,omitempty"`
	CaseTemplate string                         `json:"caseTemplate,omitempty"`
	Observables  []observable.ObservableRequest `json:"observables,omitempty"`
	Procedures   []procedure.ProcedureRequest   `json:"procedures,omitempty"`
}

type UpdateAlertRequest struct {
	Type         string         `json:"type,omitempty"`
	Source       string         `json:"source,omitempty"`
	SourceRef    string         `json:"sourceRef,omitempty"`
	ExternalLink string         `json:"externalLink,omitempty"`
	Title        string         `json:"title,omitempty"`
	Description  string         `json:"description,omitempty"`
	Severity     int            `json:"severity,omitempty"`
	Date         int64          `json:"date,omitempty"`
	LastSyncDate int64          `json:"lastSyncDate,omitempty"`
	Tags         []string       `json:"tags,omitempty"`
	AddTags      []string       `json:"addTags,omitempty"`
	RemoveTags   []string       `json:"removeTags,omitempty"`
	TLP          int            `json:"tlp,omitempty"`
	PAP          int            `json:"pap,omitempty"`
	Follow       bool           `json:"follow,omitempty"`
	CustomFields map[string]any `json:"customFields,omitempty"`
	Status       string         `json:"status,omitempty"`
	Summary      string         `json:"summary,omitempty"`
	Assignee     string         `json:"assignee,omitempty"`
}
