package api

import (
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sort"
)

var C = math.Sqrt(2.0)

const WIN_VALUE = 1
const DRAW_VALUE = 0.5
const T = 0.5001 // Control exploration with temp, T -> 0 no exploration, T->1 reflects a propablity based on visits
const MAX_CHILD_NODES = 7

//Node : MonteCarlo tree node
type Node struct {
	parent        *Node
	boardIndex    big.Int
	ChildNodes    []*Node
	unplayedMoves []int
	//Player who moved to get to this state
	playerJustMoved int64
	wins            int
	draws           int

	//Current State parameters
	VisitCount           uint      // N
	v                    float32   //value as returned by NN
	Q                    float32   //Action value ( tracks the avg of all v under this node)
	vTotal               float32   //Sum of all v under this node
	propActionChildNodes []float32 //Propablity of all child actions from this node, returned by NN

	//Edge Paramaters connecting to node
	action int     //Action taken by parent to get to this node
	p      float32 //Propablity of this action from Parent

}

func (node *Node) update(v float32) {
	node.VisitCount++
	node.vTotal += v
	//Calculate average value
	node.Q = (node.vTotal / float32(node.VisitCount))
}

func (node *Node) init(playerJustMoved int64, parent *Node, boardIndex big.Int, action int, unplayedMoves []int, propablityOfAction float32,
	propablityActionsOfChildNodes []float32, value float32) {

	node.action = action
	node.parent = parent
	node.boardIndex = boardIndex
	node.unplayedMoves = unplayedMoves
	node.playerJustMoved = playerJustMoved

	node.p = propablityOfAction
	node.propActionChildNodes = propablityActionsOfChildNodes
	node.v = value

	node.VisitCount = 0
	node.vTotal = 0
	node.Q = 0
	//fmt.Printf("---->Adding child with parent %p\n", parent)
}

func (node *Node) addChild(playerJustMoved int64, childBoardIndex big.Int, action int, childUnplayedMoves []int, propablityOfAction float32,
	propablityActionsOfChild []float32, value float32) *Node {
	var unplayedMovesAfterRemoval []int
	for _, move := range node.unplayedMoves {
		if move != action {
			unplayedMovesAfterRemoval = append(unplayedMovesAfterRemoval, move)
		}
	}
	node.unplayedMoves = unplayedMovesAfterRemoval
	var childNode Node
	childNode.init(playerJustMoved, node, childBoardIndex, action, childUnplayedMoves, propablityOfAction, propablityActionsOfChild, value)
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
	//Exploration term is high for less visited nodes
	//If child has been visited relatively less times than other children the exploration term for child goes up
	var explorationTerm = float32(C) * node.p / (1 + float32(node.VisitCount))
	//fmt.Printf("UBC = %f + %f = %f \n", node.getValue(), float32(explorationTerm), (node.getValue() + float32(explorationTerm)))
	return (node.Q + (float32(explorationTerm)))
}

/*
Calculate pi from the visits made to all child nodes
This is done on the root node of the MCTS
*/
func (node *Node) GetPi() [MAX_CHILD_NODES]float64 {
	var pi [MAX_CHILD_NODES]float64
	fmt.Println("======================")
	for _, childNode := range node.ChildNodes {
		fmt.Printf("ChildNode action: %d visit: %d, node visit: %d\n", childNode.action, childNode.VisitCount, node.VisitCount)
		pi[childNode.action] = math.Pow(float64(childNode.VisitCount), T) / math.Pow(float64(node.VisitCount), T)
	}
	return pi
}

func (node *Node) selectChildByUCT() *Node {
	var nodeWithHighestUCT *Node
	highestUCT := float32(-9999)
	for _, childNode := range node.ChildNodes {
		//fmt.Printf("Evaluating Child node for UBC: %s", childNode.ToString())
		//fmt.Printf("Parent of Child node for UBC: %s\n", childNode.parent.ToString())
		//fmt.Printf("Child Node UBC: %f\n", childNode.getUbc())
		childUbc := childNode.getUbc()
		if childUbc > highestUCT {
			//fmt.Printf("Child Node UBC: %f > highest %f\n", childUbc, highestUCT)
			nodeWithHighestUCT = childNode
			highestUCT = childUbc
		}
	}
	//fmt.Printf("----Selected Child node with highest UBC: %s\n", nodeWithHighestUCT.ToString())
	return nodeWithHighestUCT
}

func (node *Node) getUnplayedMoves() []int {
	return node.unplayedMoves
}

func (node *Node) ToString() string {
	out := fmt.Sprintf("%p :Action:%d, BoardIndex:%s len(childNodes):%d unplayedMvs:%d playerJustMvd:%s v: %f visitCount: %d Q: %f UBC=%f\n", node, node.action,
		node.boardIndex.String(), len(node.ChildNodes), node.unplayedMoves, PlayerToString(node.playerJustMoved), node.v, node.VisitCount, node.Q, node.getUbc())
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

	out := fmt.Sprintf(indentStr + "\\----" + startNode.ToString() + "\n")

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
 * TODO:
	1. Design interaction of MCTS with NN ( use twirp )
	2. Design experience generation

*/
func MonteCarloTreeSearch(game *Connect4, max_iteration int, root *Node, debug bool) *Node {

	boardIndex := game.GetBoardIndex()
	//fmt.Printf("\nMCTS root node index = %d\n", boardIndex)
	var rootNode Node
	if root == nil {
		playerWhoJustMoved := game.GetPlayerWhoJustMoved()
		unplayedMoves := game.GetValidMoves()
		root = &rootNode
		//fmt.Printf("Creating ROOT node playerJustMoved: %s, unplayedMoves %v", game.PlayerToString(playerWhoJustMoved), unplayedMoves)
		nnOut := nnForwardPass(game)
		root.init(playerWhoJustMoved, nil, boardIndex, 0, unplayedMoves, 0, nnOut.p, nnOut.value)
		root.VisitCount = 1 //Set visit count for root node as we have calculated v for this board state
		root.vTotal = nnOut.value
		if debug {
			fmt.Printf(DumpTree(root, 0))
		}
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
			//fmt.Printf("Selected node: %s\n", node.ToString())
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
			//fmt.Printf("EXP: PARENT of the child to be added %p\n", node)

			nnOut := nnForwardPass(&gameTemp)
			//fmt.Printf("EXP: Adding Child node playerJustMoved: %s, move: %d, unplayedMoves %v Value = %f\n", game.PlayerToString(playerJustMoved), move, validMoves, nnOut.value)
			tempNode := node.addChild(playerJustMoved, boardIndex, move, validMoves, node.propActionChildNodes[move], nnOut.p, nnOut.value)
			//fmt.Printf("EXP: value of child")
			//fmt.Printf("Dump parent node: %s\n", node.ToString())
			node = tempNode
			//fmt.Printf("Dump child node: %s\n", node.ToString())

		}

		//Rollout : Play the complete game from the child node just created, without making any new child nodes ( counter intutive)

		//fmt.Println("****Rollout****")
		/*
			for !gameTemp.IsGameOver() {
				playableMoves := gameTemp.GetValidMoves()
				move := playableMoves[rand.Intn(len(playableMoves))]
				gameTemp.PlayMove(move)
			}
		*/

		//Backpropagate : We should be in a terminal state when we get here
		//fmt.Println("****Backpropagate****")

		//var tempNode *Node
		rewards := []float32{0, 0} //Player who just moved
		playerJustMoveIndex := GetPlayerIndex(node.playerJustMoved)
		rewards[playerJustMoveIndex] = node.v
		rewards[(playerJustMoveIndex+1)%2] = -node.v
		for node != nil {
			playerJustMoved := node.playerJustMoved
			//fmt.Printf("BP: Player just moved: %s\n", gameTemp.PlayerToString(playerJustMoved))
			playerJustMoveIndex := GetPlayerIndex(playerJustMoved)
			//fmt.Printf("---Backpropagate---\n %s reward: %d\n", node.ToString(), rewards[playerJustMoveIndex])
			//fmt.Printf("BP: Dump current Node before: %s\n", node.ToString())
			node.update(rewards[playerJustMoveIndex])
			//fmt.Printf("BP: Dump curr   after UPDATE: %s\n", node.ToString())
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
	//fmt.Printf("Selected move for %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), root.ChildNodes[len(root.ChildNodes)-1].action)
	fmt.Printf("Pis : %v\n", root.GetPi())
	return root.ChildNodes[len(root.ChildNodes)-1]
}
