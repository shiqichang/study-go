package algorithm

import (
	"fmt"
	"testing"
)

func TestSortList(t *testing.T) {
	head := &ListNode{Val: 4}
	head.Next = &ListNode{Val: 2}
	head.Next.Next = &ListNode{Val: 1}
	head.Next.Next.Next = &ListNode{Val: 3}
	sortHead := sortList(head)
	for sortHead != nil {
		fmt.Printf("%v,", sortHead.Val)
		sortHead = sortHead.Next
	}
	fmt.Println()
}

func TestReverseList(t *testing.T) {
	head := &ListNode{Val: 1}
	head.Next = &ListNode{Val: 2}
	head.Next.Next = &ListNode{Val: 3}
	head.Next.Next.Next = &ListNode{Val: 4}
	head.Next.Next.Next.Next = &ListNode{Val: 5}
	reverseHead := reverseList(head)
	for reverseHead != nil {
		fmt.Printf("%v,", reverseHead.Val)
		reverseHead = reverseHead.Next
	}
}
