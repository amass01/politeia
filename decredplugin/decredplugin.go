// Copyright (c) 2017-2021 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package decredplugin

import (
	"encoding/json"
)

// Plugin settings, kinda doesn;t go here but for now it is fine
const (
	Version          = "1"
	ID               = "decred"
	CmdBestBlock     = "bestblock"
	CmdNewComment    = "newcomment"
	CmdCensorComment = "censorcomment"
	CmdGetComments   = "getcomments"
)

// ErrorStatusT represents decredplugin errors that result from casting a vote.
//
// These are part of the www/v1 API and must stay in until the deprecated cast
// votes route is removed.
type ErrorStatusT int

const (
	ErrorStatusInvalid          ErrorStatusT = 0
	ErrorStatusInternalError    ErrorStatusT = 1
	ErrorStatusProposalNotFound ErrorStatusT = 2
	ErrorStatusInvalidVoteBit   ErrorStatusT = 3
	ErrorStatusVoteHasEnded     ErrorStatusT = 4
	ErrorStatusDuplicateVote    ErrorStatusT = 5
	ErrorStatusIneligibleTicket ErrorStatusT = 6
	ErrorStatusLast             ErrorStatusT = 7
)

var (
	// ErrorStatus converts error status codes to human readable text.
	ErrorStatus = map[ErrorStatusT]string{
		ErrorStatusInvalid:          "invalid error status",
		ErrorStatusInternalError:    "internal error",
		ErrorStatusProposalNotFound: "proposal not found",
		ErrorStatusInvalidVoteBit:   "invalid vote bit",
		ErrorStatusVoteHasEnded:     "vote has ended",
		ErrorStatusDuplicateVote:    "duplicate vote",
		ErrorStatusIneligibleTicket: "ineligbile ticket",
	}
)

// Comment is the structure that describes the full server side content.  It
// includes server side meta-data as well. Note that the receipt is the server
// side.
type Comment struct {
	// Data generated by client
	Token     string `json:"token"`     // Censorship token
	ParentID  string `json:"parentid"`  // Parent comment ID
	Comment   string `json:"comment"`   // Comment
	Signature string `json:"signature"` // Client Signature of Token+ParentID+Comment
	PublicKey string `json:"publickey"` // Pubkey used for Signature

	// Metadata generated by decred plugin
	CommentID   string `json:"commentid"`   // Comment ID
	Receipt     string `json:"receipt"`     // Server signature of the client Signature
	Timestamp   int64  `json:"timestamp"`   // Received UNIX timestamp
	TotalVotes  uint64 `json:"totalvotes"`  // Total number of up/down votes
	ResultVotes int64  `json:"resultvotes"` // Vote score
	Censored    bool   `json:"censored"`    // Has this comment been censored
}

// EncodeComment encodes Comment into a JSON byte slice.
func EncodeComment(c Comment) ([]byte, error) {
	return json.Marshal(c)
}

// DecodeComment decodes a JSON byte slice into a Comment
func DecodeComment(payload []byte) (*Comment, error) {
	var c Comment

	err := json.Unmarshal(payload, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// NewComment sends a comment from a user to a specific proposal.  Note that
// the user is implied by the session.
type NewComment struct {
	Token     string `json:"token"`     // Censorship token
	ParentID  string `json:"parentid"`  // Parent comment ID
	Comment   string `json:"comment"`   // Comment
	Signature string `json:"signature"` // Signature of Token+ParentID+Comment
	PublicKey string `json:"publickey"` // Pubkey used for Signature
}

// EncodeNewComment encodes NewComment into a JSON byte slice.
func EncodeNewComment(nc NewComment) ([]byte, error) {
	return json.Marshal(nc)
}

// DecodeNewComment decodes a JSON byte slice into a NewComment
func DecodeNewComment(payload []byte) (*NewComment, error) {
	var nc NewComment

	err := json.Unmarshal(payload, &nc)
	if err != nil {
		return nil, err
	}

	return &nc, nil
}

// NewCommentReply returns the metadata generated by decred plugin for the new
// comment.
type NewCommentReply struct {
	CommentID string `json:"commentid"` // Comment ID
	Receipt   string `json:"receipt"`   // Server signature of the client Signature
	Timestamp int64  `json:"timestamp"` // Received UNIX timestamp
}

// EncodeNewCommentReply encodes NewCommentReply into a JSON byte slice.
func EncodeNewCommentReply(ncr NewCommentReply) ([]byte, error) {
	return json.Marshal(ncr)
}

// DecodeNewCommentReply decodes a JSON byte slice into a NewCommentReply.
func DecodeNewCommentReply(payload []byte) (*NewCommentReply, error) {
	var ncr NewCommentReply

	err := json.Unmarshal(payload, &ncr)
	if err != nil {
		return nil, err
	}

	return &ncr, nil
}

// CensorComment is a journal entry for a censored comment.  The signature and
// public key are from the admin that censored this comment.
type CensorComment struct {
	Token     string `json:"token"`     // Proposal censorship token
	CommentID string `json:"commentid"` // Comment ID
	Reason    string `json:"reason"`    // Reason comment was censored
	Signature string `json:"signature"` // Client signature of Token+CommentID+Reason
	PublicKey string `json:"publickey"` // Pubkey used for signature

	// Generated by decredplugin
	Receipt   string `json:"receipt,omitempty"`   // Server signature of client signature
	Timestamp int64  `json:"timestamp,omitempty"` // Received UNIX timestamp
}

// EncodeCensorComment encodes CensorComment into a JSON byte slice.
func EncodeCensorComment(cc CensorComment) ([]byte, error) {
	return json.Marshal(cc)
}

// DecodeCensorComment decodes a JSON byte slice into a CensorComment.
func DecodeCensorComment(payload []byte) (*CensorComment, error) {
	var cc CensorComment
	err := json.Unmarshal(payload, &cc)
	if err != nil {
		return nil, err
	}
	return &cc, nil
}

// CensorCommentReply returns the receipt for the censoring action. The
// receipt is the server side signature of CommentCensor.Signature.
type CensorCommentReply struct {
	Receipt string `json:"receipt"` // Server signature of client signature
}

// EncodeCensorCommentReply encodes CensorCommentReply into a JSON byte slice.
func EncodeCensorCommentReply(ccr CensorCommentReply) ([]byte, error) {
	return json.Marshal(ccr)
}

// DecodeCensorCommentReply decodes a JSON byte slice into a CensorCommentReply.
func DecodeCensorCommentReply(payload []byte) (*CensorCommentReply, error) {
	var ccr CensorCommentReply
	err := json.Unmarshal(payload, &ccr)
	if err != nil {
		return nil, err
	}
	return &ccr, nil
}

// GetComments retrieve all comments for a given proposal. This call returns
// the cooked comments; deleted/censored comments are not returned.
type GetComments struct {
	Token string `json:"token"` // Proposal ID
}

// EncodeGetComments encodes GetCommentsReply into a JSON byte slice.
func EncodeGetComments(gc GetComments) ([]byte, error) {
	return json.Marshal(gc)
}

// DecodeGetComments decodes a JSON byte slice into a GetComments.
func DecodeGetComments(payload []byte) (*GetComments, error) {
	var gc GetComments

	err := json.Unmarshal(payload, &gc)
	if err != nil {
		return nil, err
	}

	return &gc, nil
}

// GetCommentsReply returns the provided number of comments.
type GetCommentsReply struct {
	Comments []Comment `json:"comments"` // Comments
}

// EncodeGetCommentsReply encodes GetCommentsReply into a JSON byte slice.
func EncodeGetCommentsReply(gcr GetCommentsReply) ([]byte, error) {
	return json.Marshal(gcr)
}

// DecodeGetCommentsReply decodes a JSON byte slice into a GetCommentsReply.
func DecodeGetCommentsReply(payload []byte) (*GetCommentsReply, error) {
	var gcr GetCommentsReply

	err := json.Unmarshal(payload, &gcr)
	if err != nil {
		return nil, err
	}

	return &gcr, nil
}

// BestBlock is a command to request the best block data.
type BestBlock struct{}

// EncodeBestBlock encodes an BestBlock into a JSON byte slice.
func EncodeBestBlock(bb BestBlock) ([]byte, error) {
	return json.Marshal(bb)
}

// DecodeBestBlock decodes a JSON byte slice into a BestBlock.
func DecodeBestBlock(payload []byte) (*BestBlock, error) {
	var bb BestBlock
	err := json.Unmarshal(payload, &bb)
	if err != nil {
		return nil, err
	}
	return &bb, nil
}

// BestBlockReply is the reply to the BestBlock command.
type BestBlockReply struct {
	Height uint32 `json:"height"`
}

// EncodeBestBlockReply encodes an BestBlockReply into a JSON byte slice.
func EncodeBestBlockReply(bbr BestBlockReply) ([]byte, error) {
	return json.Marshal(bbr)
}

// DecodeBestBlockReply decodes a JSON byte slice into a BestBlockReply.
func DecodeBestBlockReply(payload []byte) (*BestBlockReply, error) {
	var bbr BestBlockReply
	err := json.Unmarshal(payload, &bbr)
	if err != nil {
		return nil, err
	}
	return &bbr, nil
}
