package mtypes

// Specify the access of a mailing list member
type AccessLevel string

// A mailing list may have one of three membership modes.
const (
	// ReadOnly specifies that nobody, including Members, may send messages to
	// the mailing list.  Messages distributed on such lists come from list
	// administrator accounts only.
	AccessLevelReadOnly = "readonly"
	// Members specifies that only those who subscribe to the mailing list may
	// send messages.
	AccessLevelMembers = "members"
	// Everyone specifies that anyone and everyone may both read and submit
	// messages to the mailing list, including non-subscribers.
	AccessLevelEveryone = "everyone"
)

// Set where replies should go
type ReplyPreference string

// Replies to a mailing list should go to one of two preferred destinations.
const (
	// List specifies that replies should be sent to the mailing list address.
	ReplyPreferenceList = "list"
	// Sender specifies that replies should be sent to the sender (FROM) address.
	ReplyPreferenceSender = "sender"
)

// A List structure provides information for a mailing list.
//
// AccessLevel may be one of ReadOnly, Members, or Everyone.
type MailingList struct {
	Address         string          `json:"address,omitempty"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	AccessLevel     AccessLevel     `json:"access_level,omitempty"`
	ReplyPreference ReplyPreference `json:"reply_preference,omitempty"`
	CreatedAt       RFC2822Time     `json:"created_at,omitempty"`
	MembersCount    int             `json:"members_count,omitempty"`
}

type ListMailingListsResponse struct {
	Items  []MailingList `json:"items"`
	Paging Paging        `json:"paging"`
}

type GetMailingListResponse struct {
	MailingList MailingList `json:"list"`
}
