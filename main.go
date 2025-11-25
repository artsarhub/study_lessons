package main

type node struct {
	V    int
	Next *node
}

func reverseList(head *node) *node {
	var prev *node
	curr := head

	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}

	return prev
}

func main() {

}
