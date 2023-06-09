/*
 * Questions: Why is 'Position win pc = -99.241974' showing negative values when the player who makes penultimate move is sure to win?
 *
 */
package api

import (
	"db"
	"fmt"
	"math"
	"math/rand"
	"mcts"
	"sort"
)

var C = 1.0 // Used CPuct optimizer on gen4 to find sweet spot //math.Sqrt(2.0) // ref:https://medium.com/oracledevs/lessons-from-alphazero-part-3-parameter-tweaking-4dceb78ed1e5  //math.Sqrt(2.0)

const T = 0.99 //0.5001 // Control exploration with temp, T -> 0 no exploration, T->1 reflects a propablity based on visits
const MAX_CHILD_NODES = 7

// Node : MonteCarlo tree node
// Node is seen from the perspective of the PLAYER who is going to move from this node ( or board state )
// Z & v are from the perspective of the player who is in this node(state) and about to make move
type Node struct {
	parent        *Node
	boardIndex    string
	ChildNodes    []*Node
	unplayedMoves []int
	//Player who moved to get to this state
	playerJustMoved int64
	wins            int
	draws           int

	//Current State parameters
	VisitCount           uint      // N
	v                    float32   //value as returned by NN , v is the value of the node from eyes of player who would be playing to get to this node
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

func (node *Node) init(playerJustMoved int64, parent *Node, boardIndex string, action int, unplayedMoves []int, propablityOfAction float32,
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

func (node *Node) addChild(playerJustMoved int64, childBoardIndex string, action int, childUnplayedMoves []int, propablityOfAction float32,
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

func (node *Node) getChildNodes() []*Node {
	return node.ChildNodes
}

func (node *Node) GetParentNode() *Node {
	return node.parent
}

// Upper bound confidence
func (node *Node) getUbc() float32 {
	//Exploration term is high for less visited nodes
	//If child has been visited relatively less times than other children the exploration term for child goes up
	//fmt.Printf("Cpuct = %f\n", C)
	var explorationTerm = float32(C) * node.p * float32(math.Sqrt(float64(node.parent.VisitCount))) / (1 + float32(node.VisitCount))
	//fmt.Printf("UBC = %f + %f = %f \n", node.getValue(), float32(explorationTerm), (node.getValue() + float32(explorationTerm)))
	return (node.Q + (float32(explorationTerm)))
}

/*
Calculate pi from the visits made to all child nodes
This is done on the root node of the MCTS

# Pi out is softmax, ie all values sum to 1

Pi[7] , where index is action
*/
func (node *Node) GetPi(printDebug bool) [MAX_CHILD_NODES]float64 {
	var pi [MAX_CHILD_NODES]float64
	if printDebug == true {
		fmt.Printf("Go Zero boardIndex = %s\n", node.boardIndex)
		fmt.Printf("==========Go Zero==ToPlay: %s==========\n", PlayerToString(node.playerJustMoved*-1))
	}
	for _, childNode := range node.ChildNodes {
		pi[childNode.action] = math.Pow(float64(childNode.VisitCount), T) / math.Pow(float64(node.VisitCount), T)
		if printDebug == true {
			fmt.Printf("ChildNode action: %d visit: %d (PC:%f), Parent node visit: %d\n", childNode.action, childNode.VisitCount, float32(100.0*childNode.VisitCount/node.VisitCount), node.VisitCount)
			fmt.Printf("Pi[%d] = %f\n", childNode.action, pi[childNode.action])
		}
	}
	if printDebug == true {
		fmt.Printf("Pi  = %v\n", pi)
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

// Propablity of all child actions from this node, returned by NN
func (node *Node) GetP() []float32 {
	return node.propActionChildNodes
}

func (node *Node) GetBoardIndex() string {
	return node.boardIndex
}

func (node *Node) GetPlayerJustMoved() int64 {
	return node.playerJustMoved
}

func (node *Node) GetQ() float32 {
	return node.Q
}

func (node *Node) GetV() float32 {
	return node.v
}

func (node *Node) ToString() string {

	if node.parent == nil {
		return "Parent Node, VisitCount:" + fmt.Sprint(node.VisitCount) + " Board Index: " + node.boardIndex
	}
	out := fmt.Sprintf("VisitCount: %d :Action:%d, BoardIndex:%s len(childNodes):%d unplayedMvs:%d playerJustMvd:%s v: %f visitCount: %d Q: %f UBC=%f\n", node.VisitCount, node.action,
		node.boardIndex, len(node.ChildNodes), node.unplayedMoves, PlayerToString(node.playerJustMoved), node.v, node.VisitCount, node.Q, node.getUbc())
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

// Return parent of node
func (node *Node) GetParent() *Node {
	return node.parent
}

/* Pick a random child from children */
func (node *Node) GetRandomChildNode() *Node {
	var numChildNodes = len(node.ChildNodes)
	fmt.Printf("Child nodes to pick from random: %d\n", numChildNodes)
	var randomChildIndex = rand.Intn(numChildNodes)

	return node.ChildNodes[randomChildIndex]
}

func CloneGame(game *Connect4) *mcts.Connect4 {
	var mctsGame mcts.Connect4

	mctsGame.SetBoard(game.board)
	mctsGame.SetPlayerMadeBadMove(game.playerMadeBadMove)
	mctsGame.SetNextPlayerToMove(game.nextPlayerToMove)
	mctsGame.SetGameOver(game.gameOver)
	mctsGame.SetReward(game.reward)

	return &mctsGame
}

var duplicateCount = 0
var database db.Database

func MonteCarloCacheInit() {
	database.CreateCacheTableMemory()
	database.CreateCacheTableFile()
	database.SyncFileDbToMemoryDb()
}

func MonteCarloCacheSyncToFile() {
	fmt.Println("Syncing in memory cache to file")
	database.SyncMemoryDbToFileDb()
}

// This variable dicates the MCTS depth when using MCTS backend ( instead of NN backend ), not to be be confused with max_iterations in func MonteCarloTreeSearch()
const MAX_MCTS_ITERATIONS = 500

var MctsIterationPercent float32 = 0.0

func MonteCarloTreeSearch(mctsTreeSearch bool, propablisticSampleOfPi bool, game *Connect4, max_iteration int, serverPort int, root *Node, debug bool) *Node {
	var cacheHits = 0
	var mctsIterations = 0
	//fmt.Printf("\nMCTSNN root node index = %s\n", boardIndex.String())
	var rootNode Node

	boardIndex := game.GetBoardIndex()
	//First time MCT is created if root == nil
	if root == nil {
		playerWhoJustMoved := game.GetPlayerWhoJustMoved()
		unplayedMoves := game.GetValidMoves()
		root = &rootNode
		//fmt.Printf("Creating ROOT node playerJustMoved: %s, unplayedMoves %v", game.PlayerToString(playerWhoJustMoved), unplayedMoves)
		//Wait for pending cache entry adds to finish
		//No wait required for first call of GetCacheEntry
		retCode, nnOut := database.GetCacheEntry(boardIndex)
		if retCode == false {
			if mctsTreeSearch {
				mctsGame := CloneGame(game)
				nnOut = mcts.MctsForwardPass(mctsGame, MAX_MCTS_ITERATIONS)
				database.InsertCacheEntry(boardIndex, nnOut)
			} else {

				//Multimove forward pass
				nnOutArray, boardIndexes := NnForwardPassMultiBoard(game, serverPort)
				database.InsertCacheEntries(boardIndexes, nnOutArray)
				if boardIndexes[MAX_CHILD_NODES] != boardIndex {
					panic("Board Indexes dont match!!")
				}
				nnOut = nnOutArray[MAX_CHILD_NODES] //Last entry is current board Index
				//single move forward pass
				//nnOut = NnForwardPass(game, serverPort)
			}
		} else {
			//boardIndex := game.GetBoardIndex()
			//fmt.Printf("Cache hit for boardIndex: %s", boardIndex.String())
		}

		root.init(playerWhoJustMoved, nil, boardIndex, 0, unplayedMoves, 0, nnOut.P, nnOut.Value)

		root.vTotal = nnOut.Value
		root.VisitCount = 0 //Visit counts are updated in update(), it comes all the way to root

		if debug {
			fmt.Printf(DumpTree(root, 0))
		}
	}
	var node *Node
	//fmt.Printf("-------->Root node = %p", root)
	for i := 0; i < max_iteration; i++ {
		fmt.Printf("\rThinking %.2f%% complete", float32(i*100)/float32(max_iteration))
		//Expose % thinking process for API
		MctsIterationPercent = float32(i*100) / float32(max_iteration)
		//fmt.Printf("\n\nMCTSNN Iteration: %d ======================================================\n", i)
		node = root
		//fmt.Printf("-------->Root node = %p\n", root)
		//Make a copy of the gamestate
		var gameTemp Connect4

		gameTemp = *game

		//Select, if all possible moves have been played, ie all possible child nodes are created
		//fmt.Println("****Select****")
		for len(node.getUnplayedMoves()) == 0 && len(node.ChildNodes) != 0 {
			//Select a child by its UCT
			node = node.selectChildByUCT()
			//fmt.Printf("Selected node: %s\n", node.ToString())
			gameTemp.PlayMove(node.action)
		}

		//#Expand
		//fmt.Printf("****Expand**** Game OVER: %t\n", gameTemp.IsGameOver())
		if len(node.getUnplayedMoves()) > 0 && gameTemp.IsGameOver() != true {
			unplayedMoves := node.getUnplayedMoves()
			//Pick a random move from unplayed Moves
			move := unplayedMoves[rand.Intn(len(unplayedMoves))]

			gameTemp.PlayMove(move)

			if gameTemp.IsGameOver() != true { 
				//With MCTS backend, MCTS forward pass does not work for completed game

				//Collect state information for new child node creation
				playerJustMoved := gameTemp.GetPlayerWhoJustMoved()
				boardIndex = gameTemp.GetBoardIndex()
				validMoves := gameTemp.GetValidMoves()
				//fmt.Printf("EXP: action: %d PARENT of the child to be added %p\n", move, node)
				mctsIterations++
				//shared.Wg.Wait()
				retCode, nnOut := database.GetCacheEntry(boardIndex)
				if retCode == false { //No cache entry found
					if mctsTreeSearch {
						mctsGame := CloneGame(&gameTemp)
						nnOut = mcts.MctsForwardPass(mctsGame, MAX_MCTS_ITERATIONS)
						database.InsertCacheEntry(boardIndex, nnOut)
						//fmt.Println("\nInserting cache entry : " + boardIndex)
					} else {
						nnOutArray, boardIndexes := NnForwardPassMultiBoard(&gameTemp, serverPort)
						if boardIndexes[MAX_CHILD_NODES] != boardIndex {
							gameTemp.DumpBoard()
							fmt.Printf("\nBoardIndexes: %v\n", boardIndexes)
							fmt.Println("Board Index evaluated = " + boardIndex + " boardIndexes[MAX_CHILD_NODES] = " + boardIndexes[MAX_CHILD_NODES] + "\n")
							panic("Board Indexes dont match!!")
						}
						//Single move forward pass
						//nnOut = NnForwardPass(&gameTemp, serverPort)
						//fmt.Println("\nInserting cache entry : " + boardIndex)

						//Multimove forward pass
						database.InsertCacheEntries(boardIndexes, nnOutArray)
						nnOut = nnOutArray[MAX_CHILD_NODES] //Last entry is current board Index
					}
				} else {
					cacheHits++
					//fmt.Println("\nCache hit for boardIndex: " + boardIndex.String())
				}

				//fmt.Printf("EXP: Adding Child node playerJustMoved: %s, move: %d, unplayedMoves %v Value = %f\n", game.PlayerToString(playerJustMoved), move, validMoves, nnOut.Value)
				tempNode := node.addChild(playerJustMoved, boardIndex, move, validMoves, node.propActionChildNodes[move], nnOut.P, nnOut.Value)
				//fmt.Printf("EXP: value of child")
				//fmt.Printf("Dump parent node: %s\n", node.ToString())
				node = tempNode
				//fmt.Printf("Dump child node: %s\n", node.ToString())
			} else { //If game is terminated create a node to capture that state
				//Collect state information for new child node creation
				playerJustMoved := gameTemp.GetPlayerWhoJustMoved()
				boardIndex = gameTemp.GetBoardIndex()
				var emptyMovesArray []int
				var zeroPropablityArray []float32
				validMoves := emptyMovesArray
				//fmt.Printf("EXP: action: %d PARENT of the child to be added %p\n", move, node)
				var value float32
				reward := gameTemp.GetReward()
				playerJustMoveIndex := GetPlayerIndex(playerJustMoved)
				//Reward is -2 or 2 for wins and 1 for draw, In MCTS it should be between -1 & 1, 0.5 for draw. Divide/2 to scale it in MCTS range
				value = float32(reward[playerJustMoveIndex]) / 2

				tempNode := node.addChild(playerJustMoved, boardIndex, move, validMoves, node.propActionChildNodes[move], zeroPropablityArray, value)
				node = tempNode
			}

		}

		//Backpropagate : We should be in a terminal state when we get here in traditional MCTS, but not in alpha algo
		//fmt.Println("****Backpropagate****")

		//TODO: Possible bug GPT analysis:
		//The backpropagation step currently updates the node based on the value of the node where the simulation ended.
		// In the AlphaGo Zero algorithm, the backpropagation step should update the node using the value from the perspective of the player who just moved.
		// You should ensure that your implementation reflects this.

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
			node = node.parent
		}
		//fmt.Printf("Root ptr: %p\n", root)
		if debug {
			fmt.Printf(DumpTree(root, 0))
		}
		/*	fmt.Printf("Temp Node %p", tempNode)
			fmt.Printf("Root Node %p", root)*/

	}

	//Now moves are sampled based on pi, so sorting is not needed but during debugs sorted order is easier to read
	sort.SliceStable(root.ChildNodes, func(i, j int) bool {
		return root.ChildNodes[i].VisitCount < root.ChildNodes[j].VisitCount
	})

	//fmt.Printf("Selected move for %s = %d\n", game.PlayerToString(game.GetPlayerToMove()), root.ChildNodes[len(root.ChildNodes)-1].action)
	const DEBUGS_TRUE = true
	pi := root.GetPi(DEBUGS_TRUE)
	//fmt.Printf("Pis : %v\n", pi)
	fmt.Printf("Position win pc = %f\n", root.Q*100)
	fmt.Printf("Cache hit PC = %f\n", 100*float32(cacheHits)/float32(mctsIterations))
	//TODO: Sampled move might not exist for a child node, As even 0 (e^0 = 1) values result in positive propablity with softmax
	//Softmax might not be a good idea for picking moves from pi, need a function which returns zero for zero pi values,
	//so that valid moves are only sampled
	var moveSampled int
	if propablisticSampleOfPi {
		moveSampled = propablisticSampleFromArray(pi)
		fmt.Printf("Propablistic sampled move from Pi: %d\n", moveSampled)
	} else {
		moveSampled = pickHighestSampleFromArray(pi)
		fmt.Printf("Highest visited move from Pi: %d\n", moveSampled)
	}

	for _, childNode := range root.ChildNodes {
		if childNode.GetAction() == moveSampled {
			return childNode
		}
	}
	panic(" Could not find sampled child node")
	// Return the child with most visits
	//return root.ChildNodes[len(root.ChildNodes)-1]
}

func sum(array [MAX_CHILD_NODES]float64) float64 {
	result := 0.00
	for _, v := range array {
		result += v
	}
	return result
}

func pickHighestSampleFromArray(in [MAX_CHILD_NODES]float64) int {
	highestIndex := 0
	highestValue := 0.0
	for index, value := range in {
		if value > highestValue {
			highestIndex = index
			highestValue = value
		}
	}
	return highestIndex
}

/*
 *  Sample from array based on the value, zero values of array never get selected
 */

func propablisticSampleFromArray(in [MAX_CHILD_NODES]float64) int {

	randFloat := RandomNumGenerator.Float64()
	pick := randFloat * sum(in)
	//fmt.Printf("Pick = %f\n", pick)
	sum := 0.0
	for index, value := range in {
		sum += value
		if pick < sum {
			return index
		}
	}
	//The pick is the last node of the list
	return MAX_CHILD_NODES - 1

}

/*
//
//  Pick a propablistic index from a softmax out. The propablity of selection of an index is propotional to the value at the index.
//
func sampleFromSoftmax(softMaxOut [MAX_CHILD_NODES]float64) int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	var result int
	var sum float64 = 0

	randFloat := r1.Float64() //rand.Float64(r1)
	fmt.Println(randFloat)
	for i := 0; i < len(softMaxOut); i++ {
		sum += softMaxOut[i]
		if randFloat <= sum {
			//fmt.Printf("Selection = %d\n", i)
			result = i
			break
		}
	}
	return result
}

func softMax(x [MAX_CHILD_NODES]float64) [MAX_CHILD_NODES]float64 {
	var denominator float64 = 0
	for _, value := range x {
		denominator += math.Exp(value)
	}

	var a [MAX_CHILD_NODES]float64

	for index, value := range x {
		a[index] = math.Exp(value) / denominator
	}

	return a
}


/*
func checkDuplicateBoardIndex(boardIndex big.Int) {
	var found bool = false
	for _, bIndexCache := range cache {
		if bIndexCache.Cmp(&boardIndex) == 0 {
			found = true
			break
		}
	}

	if !found {
		cache = append(cache, boardIndex)
	} else {
		duplicateCount++
		fmt.Printf("%s Duplicate entry!!!!!!!!!!!!!!!!Count = %d\n", boardIndex.String(), duplicateCount)
	}
}
*/
