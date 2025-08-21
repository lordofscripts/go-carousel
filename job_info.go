/* -----------------------------------------------------------------
 *					L o r d  O f   S c r i p t s (tm)
 *				  Copyright (C)2025 Dídimo Grimaldo T.
 *							   workspace.json
 * - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
 *
 *-----------------------------------------------------------------*/
package carousel

/* ----------------------------------------------------------------
 *				P u b l i c		T y p e s
 *-----------------------------------------------------------------*/

type JobInfo struct {
	Id        uint
	TimeStamp string
	Title     string
}

type JobInfoSlice []JobInfo

/* ----------------------------------------------------------------
 *				P u b l i c		M e t h o d s
 *-----------------------------------------------------------------*/

// Implement the sort.Interface based on the TimeStamp field
// Usage: sort.Sort(JobInfoSlice(array))
func (ji JobInfoSlice) Len() int           { return len(ji) }
func (ji JobInfoSlice) Swap(i, j int)      { ji[i], ji[j] = ji[j], ji[i] }
func (ji JobInfoSlice) Less(i, j int) bool { return ji[i].TimeStamp < ji[j].TimeStamp } // ascending
