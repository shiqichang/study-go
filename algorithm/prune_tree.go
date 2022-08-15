package algorithm

/*
二叉树剪枝：
给定一个二叉树 根节点root，树的每个节点的值要么是 0，要么是 1。请剪除该二叉树中所有节点的值为 0 的子树。
节点 node 的子树为node 本身，以及所有 node的后代。

示例 1:
输入: [1,null,0,0,1]
输出: [1,null,0,null,1]
解释:

示例 2:
输入: [1,0,1,0,0,0,1]
输出: [1,null,1,null,1]
 */

func pruneTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	root.Left = pruneTree(root.Left)
	root.Right = pruneTree(root.Right)
	if root.Left == nil && root.Right == nil && root.Val == 0 {
		return nil
	}
	return root
}
