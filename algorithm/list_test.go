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

func TestAddTwoNumbers(t *testing.T) {
	l1 := &ListNode{Val: 7}
	l1.Next = &ListNode{Val: 1}
	l1.Next.Next = &ListNode{Val: 6}

	l2 := &ListNode{Val: 5}
	l2.Next = &ListNode{Val: 9}
	l2.Next.Next = &ListNode{Val: 3}

	res := addTwoNumbers(l1, l2)
	for res != nil {
		fmt.Printf("%v,", res.Val)
		res = res.Next
	}
}
