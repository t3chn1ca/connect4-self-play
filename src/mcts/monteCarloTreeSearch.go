package mcts

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sort"
)

var C = math.Sqrt(2.0) //2.00 // TODO: revert back to sqrt 2 math.Sqrt(2.0)

const T = 0.99 // 0.5001 // Control exploration with temp, T -> 0 no exploration, T->1 reflects a propablity based on visits

const WIN_VALUE = 1
const DRAW_VALUE = 0.5
const MAX_CHILD_NODES = 7

//Node : MonteCarlo tree node
type Node struct {
	action        int //Action taken by parent to get to this node
	parent        *Node
	boardIndex    big.Int
	ChildNodes    []*Node
	unplayedMoves []int
	//Player who moved to get to this state
	playerJustMoved int64
	VisitCount      int
	wins            int
	draws           int
}

func (node *Node) init(playerJustMoved int64, parent *Node, boardIndex big.Int, action int, unplayedMoves []int) {
	node.action = action
	node.parent = parent
	node.boardIndex = boardIndex
	node.unplayedMoves = unplayedMoves
	node.playerJustMoved = playerJustMoved
	//fmt.Printf("---->Adding child with parent %p\n", parent)
}

func (node *Node) addChild(playerJustMoved int64, childBoardIndex big.Int, action int, childUnplayedMoves []int) *Node {
	var unplayedMovesAfterRemoval []int
	for _, move := range node.unplayedMoves {
		if move != action {
			unplayedMovesAfterRemoval = append(unplayedMovesAfterRemoval, move)
		}
	}
	node.unplayedMoves = unplayedMovesAfterRemoval
	var childNode Node
	childNode.init(playerJustMoved, node, childBoardIndex, action, childUnplayedMoves)
	node.ChildNodes = append(node.ChildNodes, &childNode)
	return &childNode
}

func (node *Node) getValue() float32 {
	return ((float32(node.wins) * WIN_VALUE) + (float32(node.draws) * DRAW_VALUE)) / float32(node.VisitCount)
}

func (node *Node) getChildNodes() []*Node {
	return node.ChildNodes
}

// Upper bound confidence
func (node *Node) getUbc() float32 {
	//Exploration term is high for less visited nodes ( Exploration term has no  relation to wins/losses)
	//If child has been visited relatively less times than other children the exploration term for child goes up
	var explorationTerm = C * math.Sqrt(math.Log(float64(node.parent.VisitCount))/float64(node.VisitCount))
	//fmt.Printf("UBC = %f + %f = %f \n", node.getValue(), float32(explorationTerm), (node.getValue() + float32(explorationTerm)))
	return (node.getValue() + (float32(explorationTerm)))
}

func (node *Node) selectChildByUCT() *Node {
	var nodeWithHighestUCT *Node
	highestUCT := float32(-9999)
	for _, childNode := range node.ChildNodes {
		//fmt.Printf("Evaluating Child node for UBC: %s", childNode.toString())
		//fmt.Printf("Parent of Child node for UBC: %s\n", childNode.parent.toString())
		//fmt.Printf("Child Node UBC: %f\n", childNode.getUbc())
		childUbc := childNode.getUbc()
		if childUbc > highestUCT {
			//fmt.Printf("Child Node UBC: %f > highest %f\n", childUbc, highestUCT)
			nodeWithHighestUCT = childNode
			highestUCT = childUbc
		}
	}
	//fmt.Printf("----Selected Child node with highest UBC: %s\n", nodeWithHighestUCT.toString())
	return nodeWithHighestUCT
}

func (node *Node) getUnplayedMoves() []int {
	return node.unplayedMoves
}

func (node *Node) update(reward int) {
	node.VisitCount++
	if reward == 2 {
		node.wins++
	}
	if reward == 1 {
		node.draws++
	}
	if reward == -2 {
		//Dont do anything for loss
		node.wins--
	}
}

func (node *Node) toString() string {
	out := fmt.Sprintf("%p :Action:%d, BoardIndex:%s len(childNodes):%d unplayedMvs:%d playerJustMvd:%d W+D/V: %d/%d =%f \n", node, node.action,
		node.boardIndex.String(), len(node.ChildNodes), node.unplayedMoves, node.playerJustMoved, node.wins+node.draws, node.VisitCount, float32(node.wins+node.draws)/float32(node.VisitCount))
	/*  Print parent and child
	out += fmt.Sprintf("Parent: %p\n", node.parent)
	for i := 0; i < len(node.ChildNodes); i++ {
		out += fmt.Sprintf("Child:%p,", &node.ChildNodes[i])
	}

	out += "\n"
	*/
	return string(out)
}

func DumpTree(startNode *Node, indentCount int) string {

	const PER_LEVEL_INDENT = 7
	indentStr := ""
	for i := 0; i < indentCount; i++ {
		indentStr += " "
	}

	out := fmt.Sprintf(indentStr + "\\----" + startNode.toString() + "\n")

	indentCount += PER_LEVEL_INDENT

	for _, child := range startNode.ChildNodes {
		out += DumpTree(child, indentCount)
	}

	return out

}

func (node *Node) GetAction() int {
	return node.action
}

/*
Calculate pi from the visits made to all child nodes
This is done on the root node of the MCTS

Pi out is softmax, ie all values sum to 1

Pi[7] , where index is action
*/
func (node *Node) GetPi() [MAX_CHILD_NODES]float32 {
	var pi [MAX_CHILD_NODES]float32
	fmt.Println("==========MCTS============")
	for _, childNode := range node.ChildNodes {
		fmt.Printf("ChildNode action: %d visit: %d, Parent node visit: %d\n", childNode.action, childNode.VisitCount, node.VisitCount)
		pi[childNode.action] = float32(math.Pow(float64(childNode.VisitCount), T) / math.Pow(float64(node.VisitCount), T))
		fmt.Printf("Pi[%d] = %f\n", childNode.action, pi[childNode.action])
	}
	//fmt.Printf("Pi  = %v\n", pi)
	return pi
}

type NNOut struct {
	Value float32
	P     []float32
}

const MAX_MCTS_ITERATIONS = 500

func MctsForwardPass(game *Connect4) NNOut {
	selectedNode := MonteCarloTreeSearch(game, MAX_MCTS_ITERATIONS, nil, false)

	var nnOut NNOut

	parentNode := selectedNode.parent
	nnOut.Value = parentNode.getValue()
	boardIndex := game.GetBoardIndex()
	fmt.Printf("MCTS boardIndex = %s\n", boardIndex.String())
	fmt.Printf("MCTS Value: %f\n", parentNode.getValue())
	pi := parentNode.GetPi()
	nnOut.P = pi[:]

	return nnOut

}
func MonteCarloTreeSearch(game *Connect4, max_iteration int, root *Node, debug bool) *Node {

	boardIndex := game.GetBoardIndex()
	//fmt.Printf("\nMCTS root node index = %s\n", boardIndex.String())
	var rootNode Node
	if root == nil {
		playerWhoJustMoved := game.GetPlayerWhoJustMoved()
		unplayedMoves := game.GetValidMoves()
		root = &rootNode
		//fmt.Printf("Creating ROOT node playerJustMoved: %s, unplayedMoves %v", game.PlayerToString(playerWhoJustMoved), unplayedMoves)
		root.init(playerWhoJustMoved, nil, boardIndex, 0, unplayedMoves)
	}
	var node *Node

	//fmt.Printf("-------->Root node = %p", root)
	for i := 0; i < max_iteration; i++ {
		//fmt.Printf("\n\nIteration: %d ======================================================\n", i)
		node = root
		//fmt.Printf("-------->Root node = %p\n", root)
		//Make a copy of the gamestate
		var gameTemp Connect4

		gameTemp = *game

		//Select, if all possible moves have been played, ie all possible child nodes are created
		//fmt.Println("****Select****")
		for len(node.getUnplayedMoves()) == 0 && len(node.ChildNodes) != 0 {
			node = node.selectChildByUCT()
			//fmt.Printf("Selected node: %s\n", node.toString())
			gameTemp.PlayMove(node.action)
		}

		//#Expand
		//fmt.Printf("****Expand**** Game OVER: %t\n", gameTemp.IsGameOver())
		if len(node.getUnplayedMoves()) > 0 && gameTemp.IsGameOver() != true {
			unplayedMoves := node.getUnplayedMoves()

			move := unplayedMoves[rand.Intn(len(unplayedMoves))]
			//DEBUG: Remove and revert to random
			/*
				minMove := -1
				for i, e := range unplayedMoves {
					if i == 0 || e < minMove {
						minMove = e
					}
				}
				move := minMove*/
			gameTemp.PlayMove(move)

			//Collect state information for new child node creation
			playerJustMoved := gameTemp.GetPlayerWhoJustMoved()
			boardIndex = gameTemp.GetBoardIndex()
			validMoves := gameTemp.GetValidMoves()
			//fmt.Printf("PARENT of the child to be added %p\n", node)
			//fmt.Printf("Adding Child node playerJustMoved: %s, move: %d, unplayedMoves %v\n", game.PlayerToString(playerJustMoved), move, validMoves)

			tempNode := node.addChild(playerJustMoved, boardIndex, move, validMoves)
			//fmt.Printf("Dump parent node: %s\n", node.toString())
			node = tempNode
			//fmt.Printf("Dump child node: %s\n", node.toString())

		}

		//Rollout : Play the complete game from the child node just created, without making any new child nodes ( counter intutive)
		//fmt.Println("****Rollout****")
		for !gameTemp.IsGameOver() {
			playableMoves := gameTemp.GetValidMoves()
			move := playableMoves[rand.Intn(len(playableMoves))]
			//DEBUG: Revert back to random
			/*
				minMove := -1
				for i, e := range playableMoves {
					if i == 0 || e < minMove {
						minMove = e
					}
				}
			*/
			gameTemp.PlayMove(move)
		}

		//Backpropagate : We should be in a terminal state when we get here
		//fmt.Println("****Backpropagate****")

		//var tempNode *Node
		rewards := gameTemp.GetReward()
		for node != nil {
			playerJustMoved := node.playerJustMoved
			//fmt.Printf("Player just moved: %s\n", gameTemp.PlayerToString(playerJustMoved))
			playerJustMoveIndex := GetPlayerIndex(playerJustMoved)
			//fmt.Printf("---Backpropagate---\n %s reward: %d\n", node.toString(), rewards[playerJustMoveIndex])
			//fmt.Printf("Dump current Node before: %s\n", node.toString())
			node.update(rewards[playerJustMoveIndex])
			//fmt.Printf("Dump current Node after -update-: %s\n", node.toString())
			//fmt.Printf("--------------------------------\n")
			//tempNode = node
			node = node.parent
		}
		//fmt.Printf("Root ptr: %p\n", root)
		if debug {
			fmt.Printf(DumpTree(root, 0))
		}
		/*	fmt.Printf("Temp Node %p", tempNode)
			fmt.Printf("Root Node %p", root)
			fmt.Println(root.TreeToString(0))*/
	}

	sort.SliceStable(root.ChildNodes, func(i, j int) bool {
		return root.ChildNodes[i].VisitCount < root.ChildNodes[j].VisitCount
	})
	/*
		for _, node := range root.ChildNodes {
			fmt.Printf("%v\n", node.toString())
		}
	*/
	boardIndex = game.GetBoardIndex()
	fmt.Printf("\n BoardIndex %s \n", boardIndex.String())
	//fmt.Printf("\nBoardIndex: %s Selected move for %s = %d\n", boardIndex.String(), game.PlayerToString(game.GetPlayerToMove()), root.ChildNodes[len(root.ChildNodes)-1].action)
	fmt.Printf("Len childNodes = %d\n", len(root.ChildNodes))
	return root.ChildNodes[len(root.ChildNodes)-1]
}
