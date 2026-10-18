// Package parser extracts links from downloaded HTML and CSS documents.
//
// It looks at a[href], link[href] and img[wget] tags in HTML and at
// url(...) references in CSS, and resolves relative links against the
// document URL. (Phase 5)
package parser
