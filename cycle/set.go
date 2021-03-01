package cycle

// Set represents a group of values using a cyclically linked list
type Set struct {
	head *node
}
type node struct {
	v          Val
	next, prev *node
}

// Val is a member of a Set.
type Val = int

// Add inserts a new value into the set. Duplicates are allowed.
func (s *Set) Add(v Val) {
	n := &node{v: v}
	if s.head == nil {
		n.next = n
		n.prev = n
		s.head = n
	} else {
		n.next = s.head.next
		s.head.next = n
		n.next.prev = n
		n.prev = s.head
	}
}

func (s *Set) toSlice() []Val {
	if s.head == nil {
		return nil
	}
	vs := []Val{s.head.v}
	curr := s.head.next
	for curr != s.head {
		vs = append(vs, curr.v)
		curr = curr.next
	}
	return vs
}

// Merge will cause the contents of both sets to be equal to their union. Note
// that because Add allows duplicates, the size of the merged set is exactly the
// sum of the two sets' sizes.
func (s *Set) Merge(s2 *Set) {
	// Trivial cases if one or both head are nil
	if s.head == nil && s2.head == nil {
		// This is the only case where subsequent Adds aren't reflected in s2.
		return
	}
	if s.head == nil && s2.head != nil {
		s.head = s2.head
		return
	}
	if s2.head == nil && s.head != nil {
		s2.head = s.head
		return
	}
	// Merge case: "twist" any two nodes from the two sets, such that node
	// traversal from anywhere will visit both sets.
	s.head.next, s2.head.next = s2.head.next, s.head.next
	s.head.next.prev = s.head
	s2.head.next.prev = s2.head
}
