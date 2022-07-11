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
