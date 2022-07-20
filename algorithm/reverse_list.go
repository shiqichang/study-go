package algorithm

/*
给定单链表的头节点 head ，请反转链表，并返回反转后的链表的头节点。

示例 1：
输入：head = [1,2,3,4,5]
输出：[5,4,3,2,1]
 */

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	curr := head
	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}

	return prev
}
