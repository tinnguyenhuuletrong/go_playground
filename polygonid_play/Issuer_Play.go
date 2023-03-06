package polygonid_play

// https://github.com/0xPolygonID/tutorial-examples/blob/main/issuer-protocol/main.go

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/iden3/go-circuits"
	core "github.com/iden3/go-iden3-core"
	"github.com/iden3/go-iden3-crypto/babyjub"
	"github.com/iden3/go-iden3-crypto/poseidon"
	merkletree "github.com/iden3/go-merkletree-sql/v2"
	"github.com/iden3/go-merkletree-sql/v2/db/memory"
)

func Issuer_Protocol() {
	// 1. BabyJubJub key
	fmt.Println("--- BabyJubJub key ---")

	// generate babyJubjub private key randomly
	babyJubjubPrivKey := babyjub.NewRandPrivKey()

	// generate public key from private key
	babyJubjubPubKey := babyJubjubPrivKey.Public()

	// print public key
	fmt.Println(babyJubjubPubKey)

	// 2. Sparse Merkle Tree
	fmt.Println("--- Sparse Merkle Tree ---")

	ctx := context.Background()

	// Tree storage
	store := memory.NewMemoryStorage()

	// Generate a new MerkleTree with 32 levels
	mt, _ := merkletree.NewMerkleTree(ctx, store, 32)

	// Add a leaf to the tree with index 1 and value 10
	index1 := big.NewInt(1)
	value1 := big.NewInt(10)
	mt.Add(ctx, index1, value1)

	// Add another leaf to the tree
	index2 := big.NewInt(2)
	value2 := big.NewInt(15)
	mt.Add(ctx, index2, value2)

	// Proof of membership of a leaf with index 1
	proofExist, value, _ := mt.GenerateProof(ctx, index1, mt.Root())

	var jsonStr, _ = proofExist.MarshalJSON()
	fmt.Println("Proof of index 1: ", string(jsonStr))
	fmt.Println("Proof of membership:", proofExist.Existence, "Value corresponding to the queried index:", value)

	// Proof of non-membership of a leaf with index 4
	proofNotExist, _, _ := mt.GenerateProof(ctx, big.NewInt(4), mt.Root())
	jsonStr, _ = proofNotExist.MarshalJSON()
	fmt.Println("Proof of index 4: ", string(jsonStr))
	fmt.Println("Proof of membership:", proofNotExist.Existence)

	fmt.Println("Dump tree to graphviz: ", "./tmp/three.graphviz")
	f, _ := os.Create("./tmp/three.graphviz")
	mt.GraphViz(ctx, f, mt.Root())
	f.Close()

	// 3.1. Create a Generic Claim
	fmt.Println("--- Create a Generic Claim ---")

	// set claim expriation date to 2361-03-22T00:44:48+05:30
	t := time.Date(2361, 3, 22, 0, 44, 48, 0, time.UTC)

	// set schema
	ageSchema, _ := core.NewSchemaHashFromHex("2e2d1c11ad3e500de68d7ce16a0a559e")

	// define data slots
	birthday := big.NewInt(19960424)
	documentType := big.NewInt(1)

	// set revocation nonce
	revocationNonce := uint64(1909830690)

	// set ID of the claim subject
	subjectId, _ := core.IDFromString("113TCVw5KMeMp99Qdvub9Mssfz7krL9jWNvbdB7Fd2")

	// create claim
	claim, _ := core.NewClaim(ageSchema, core.WithExpirationDate(t), core.WithRevocationNonce(revocationNonce), core.WithIndexID(subjectId), core.WithIndexDataInts(birthday, documentType))

	// transform claim from bytes array to json
	claimToMarshal, _ := json.Marshal(claim)
	fmt.Println(string(claimToMarshal))

	fmt.Println("Claim HiHv: ")
	fmt.Println(claim.HiHv())

	// ["3613283249068442770038516118105710406958","0","19960424","1","227737944108667786680629310498","0","0","0"]
	// Index:
	// {
	// "3613283249068442770038516118105710406958", // Claim Schema hash
	// "86645363564555144061174553487309804257148595648980197130928167920533372928", // ID Subject of the claim
	// "19960424", // First index data slot stores the date of birth
	// "1"  //  Second index data slot stores the document type
	// }

	// Value:
	// {
	// "227737944108667786680629310498", // Revocation nonce
	// "0",
	// "0", // first value data slot
	// "0"  // second value data slot
	// }

	// 3.2. Create Auth Claim
	fmt.Println("--- Create Auth Claim ---")

	// Add revocation nonce. Used to invalidate the claim. This may be a random number in the real implementation.
	revNonce := uint64(1)

	authClaim, _ := core.NewClaim(core.AuthSchemaHash,
		core.WithIndexDataInts(babyJubjubPubKey.X, babyJubjubPubKey.Y),
		core.WithRevocationNonce(revNonce))

	authClaimToMarshal, _ := json.Marshal(authClaim)

	fmt.Println(string(authClaimToMarshal))

	// 4.1. Generate identity trees, Retrieve identity state, Retrieve Identifier (ID)
	// https://docs.iden3.io/getting-started/identity/identity-state/
	fmt.Println("--- Generate identity trees ---")

	// Create empty Claims tree
	clt, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	// Create empty Revocation tree
	ret, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	// Create empty Roots tree
	rot, _ := merkletree.NewMerkleTree(ctx, memory.NewMemoryStorage(), 32)

	// Get the Index of the claim and the Value of the authClaim
	hIndex, hValue, _ := authClaim.HiHv()

	// add auth claim to claims tree with value hValue at index hIndex
	clt.Add(ctx, hIndex, hValue)

	// print the roots
	fmt.Println(clt.Root().BigInt(), ret.Root().BigInt(), rot.Root().BigInt())

	state, _ := merkletree.HashElems(
		clt.Root().BigInt(),
		ret.Root().BigInt(),
		rot.Root().BigInt())

	fmt.Println("Identity State:", state.Hex())
	id, _ := core.IdGenesisFromIdenState(core.TypeDefault, state.BigInt())
	fmt.Println("ID:", id)

	// 5. Issuing Claim by Signature
	fmt.Println("--- Issuing Claim by Signature ---")
	//https://docs.iden3.io/getting-started/issue-claim-overview/

	// Retrieve indexHash and valueHash of the claim
	claimIndex, claimValue := claim.RawSlots()
	indexHash, _ := poseidon.Hash(core.ElemBytesToInts(claimIndex[:]))
	valueHash, _ := poseidon.Hash(core.ElemBytesToInts(claimValue[:]))

	// Poseidon Hash the indexHash and the valueHash together to get the claimHash
	claimHash, _ := merkletree.HashElems(indexHash, valueHash)

	// Sign the claimHash with the private key of the issuer
	claimSignature := babyJubjubPrivKey.SignPoseidon(claimHash.BigInt())
	fmt.Println("Claim Signature:", claimSignature)

	// 6. Issuing Claim by adding it to the Merkle Tree
	fmt.Println("--- Issuing Claim by adding it to the Merkle Tree ---")

	// GENESIS STATE:

	// 1. Generate Merkle Tree Proof for authClaim at Genesis State
	authMTPProof, _, _ := clt.GenerateProof(ctx, hIndex, clt.Root())

	// 2. Generate the Non-Revocation Merkle tree proof for the authClaim at Genesis State
	authNonRevMTPProof, _, _ := ret.GenerateProof(ctx, new(big.Int).SetUint64(revNonce), ret.Root())

	// Snapshot of the Genesis State
	genesisTreeState := circuits.TreeState{
		State:          state,
		ClaimsRoot:     clt.Root(),
		RevocationRoot: ret.Root(),
		RootOfRoots:    rot.Root(),
	}
	// STATE 1:

	// Before updating the claims tree, add the claims tree root at Genesis state to the Roots tree.
	rot.Add(ctx, clt.Root().BigInt(), big.NewInt(0))

	// Create a new random claim
	schemaHex := hex.EncodeToString([]byte("myAge_test_claim"))
	schema, _ := core.NewSchemaHashFromHex(schemaHex)

	code := big.NewInt(51)

	newClaim, _ := core.NewClaim(schema, core.WithIndexDataInts(code, nil))

	// Get hash Index and hash Value of the new claim
	hi, hv, _ := newClaim.HiHv()

	// Add claim to the Claims tree
	clt.Add(ctx, hi, hv)

	// Fetch the new Identity State
	newState, _ := merkletree.HashElems(
		clt.Root().BigInt(),
		ret.Root().BigInt(),
		rot.Root().BigInt())

	// Snapshot of the new tree State
	newTreeState := circuits.TreeState{
		State:          newState,
		ClaimsRoot:     clt.Root(),
		RevocationRoot: ret.Root(),
		RootOfRoots:    rot.Root(),
	}

	// Sign a message (hash of the genesis state + the new state) using your private key
	hashOldAndNewStates, _ := poseidon.Hash([]*big.Int{state.BigInt(), newState.BigInt()})

	signature := babyJubjubPrivKey.SignPoseidon(hashOldAndNewStates)

	authClaimNewStateIncMtp, _, _ := clt.GenerateProof(ctx, hIndex, newTreeState.ClaimsRoot)

	// Generate state transition inputs
	stateTransitionInputs := circuits.StateTransitionInputs{
		ID:                      id,
		OldTreeState:            genesisTreeState,
		NewTreeState:            newTreeState,
		IsOldStateGenesis:       true,
		AuthClaim:               authClaim,
		AuthClaimIncMtp:         authMTPProof,
		AuthClaimNonRevMtp:      authNonRevMTPProof,
		AuthClaimNewStateIncMtp: authClaimNewStateIncMtp,
		Signature:               signature,
	}

	// Perform marshalling of the state transition inputs
	inputBytes, _ := stateTransitionInputs.InputsMarshal()

	fmt.Println("SmartContract input:", string(inputBytes))

}
