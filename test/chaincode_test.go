package test

import (
	"testing"
	"time"

	"github.com/GeorgeGogos/SCENE-Chaincode/payload"

	"math/rand"

	scene "github.com/GeorgeGogos/SCENE-Chaincode"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testcc "github.com/hyperledger-labs/cckit/testing"
	expectcc "github.com/hyperledger-labs/cckit/testing/expect"
)

func TestChaincode(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SCENE Suite")
}

func randomOrgs() (string, string, string, string) {
	rand.Seed(time.Now().Unix())
	org1 := uuid.New().String()
	org2 := uuid.New().String()
	org3 := uuid.New().String()
	org4 := uuid.New().String()
	return org1, org2, org3, org4
}

func randomItems() (string, string, string, string) {
	rand.Seed(time.Now().Unix())
	item1 := uuid.New().String()
	item2 := uuid.New().String()
	item3 := uuid.New().String()
	item4 := uuid.New().String()
	return item1, item2, item3, item4
}

func randomUnits() (string, string, string, string) {
	rand.Seed(time.Now().Unix())
	unit1 := uuid.New().String()
	unit2 := uuid.New().String()
	unit3 := uuid.New().String()
	unit4 := uuid.New().String()
	return unit1, unit2, unit3, unit4
}

func randomType() (string, string, string, string) {
	Object_Type = []string{"Service", "Device", "Marketplace"}
	type1 := Object_Type[rand.Intn(len(Object_Type))]
	type2 := Object_Type[rand.Intn(len(Object_Type))]
	type3 := Object_Type[rand.Intn(len(Object_Type))]
	type4 := Object_Type[rand.Intn(len(Object_Type))]
	return type1, type2, type3, type4
}

var (
	userCN                     = "ggogos@iti.gr"
	org1, org2, org3, org4     = randomOrgs()
	item1, item2, item3, item4 = randomItems()
	type1, type2, type3, type4 = randomType()
	userID1, _                 = GenerateCertIdentity(`SomeMSP`, userCN, org1)
	userID2, _                 = GenerateCertIdentity(`SomeMSP`, userCN, org2)
	userID3, _                 = GenerateCertIdentity(`SomeMSP`, userCN, org3)
	userID4, _                 = GenerateCertIdentity(`SomeMSP`, userCN, org4)
	maliciousUserID, _         = GenerateCertIdentity(`SomeMSP`, userCN, "maliciousOrg")
	t                          = true
	f                          = false
	Object_Type                = []string{"Service", "Device", "Marketplace"}
	Contract_Type              = []string{"Private", "Community"}
	Contract_Status            = []string{"Pending", "Accepted", "Rejected"}
)

var _ = Describe(`Chaincode`, func() {

	//Create chaincode mock
	cc := testcc.NewMockStub(`scene_chaincode`, scene.NewCC())

	BeforeSuite(func() {
		// init chaincode
		expectcc.ResponseOk(cc.From(userID1).Init()) // init chaincode from authority
	})

	Describe("Contract functions (happy path)", func() {
		It("ProposeContract, expected to be succeed", func() {
			//invoke chaincode method from authority actor

			ccResponse := (cc.From(userID1).Invoke(`ProposeContract`, &payload.ContractPayload{
				ContractId:   "ID-01",
				ContractType: Contract_Type[rand.Intn(len(Contract_Type))],
				Orgs:         []string{org1, org2},
				Items: []payload.Item{{
					Enabled:    &t,
					Write:      &t,
					ObjectId:   item1,
					OrgId:      org1,
					ObjectType: type1,
				},
					{
						Enabled:    &t,
						Write:      &f,
						ObjectId:   item2,
						OrgId:      org2,
						ObjectType: type2,
					}},
			}))
			expectcc.ResponseOk(ccResponse)

			ccResponse = (cc.From(userID2).Invoke(`ProposeContract`, &payload.ContractPayload{
				ContractId:   "ID-02",
				ContractType: Contract_Type[rand.Intn(len(Contract_Type))],
				Orgs:         []string{org2, org3},
				Items: []payload.Item{{
					Enabled:    &t,
					Write:      &t,
					ObjectId:   item2,
					OrgId:      org2,
					ObjectType: type2,
				},
					{
						Enabled:    &t,
						Write:      &f,
						ObjectId:   item3,
						OrgId:      org3,
						ObjectType: type3,
					}},
			}))
			expectcc.ResponseOk(ccResponse)

			ccResponse = (cc.From(userID3).Invoke(`ProposeContract`, &payload.ContractPayload{
				ContractId:   "ID-03",
				ContractType: Contract_Type[rand.Intn(len(Contract_Type))],
				Orgs:         []string{org3, org4},
				Items: []payload.Item{{
					Enabled:    &t,
					Write:      &t,
					ObjectId:   item3,
					OrgId:      org3,
					ObjectType: type3,
				},
					{
						Enabled:    &t,
						Write:      &f,
						ObjectId:   item4,
						OrgId:      org4,
						ObjectType: type4,
					}},
			}))
			expectcc.ResponseOk(ccResponse)

			ccResponse = (cc.From(userID4).Invoke(`ProposeContract`, &payload.ContractPayload{
				ContractId:   "ID-04",
				ContractType: Contract_Type[rand.Intn(len(Contract_Type))],
				Orgs:         []string{org4, org1},
				Items: []payload.Item{{
					Enabled:    &t,
					Write:      &t,
					ObjectId:   item4,
					OrgId:      org4,
					ObjectType: type4,
				},
					{
						Enabled:    &t,
						Write:      &f,
						ObjectId:   item1,
						OrgId:      org1,
						ObjectType: type1,
					}},
			}))
			expectcc.ResponseOk(ccResponse)

		})

		It("AcceptContract from Org included in Contract, expected to succeed", func() {

			testID := "ID-01"
			ccResponse := (cc.From(userID2).Invoke(`AcceptContract`, testID))
			expectcc.ResponseOk(ccResponse)
			testID = "ID-04"
			ccResponse = (cc.From(userID1).Invoke(`AcceptContract`, testID))
			expectcc.ResponseOk(ccResponse)
			testID = "ID-03"
			ccResponse = (cc.From(userID4).Invoke(`AcceptContract`, testID))
			expectcc.ResponseOk(ccResponse)
		})

		It("RejectContract from Org included in Contract, expected to succeed", func() {

			testID := "ID-02"
			ccResponse := (cc.From(userID3).Invoke(`RejectContract`, testID))
			expectcc.ResponseOk(ccResponse)
		})

		It("DeleteContract from Org included in Contract, expected to succeed", func() {
			testID := "ID-03"
			ccResponse := (cc.From(userID3).Invoke(`DissolveContract`, testID))
			expectcc.ResponseOk(ccResponse)
		})
		It("GetContractByID from Org included in Contract, expected to succeed", func() {
			testID := "ID-01"
			ccResponse := (cc.From(userID1).Query(`GetContractByID`, testID))
			expectcc.ResponseOk(ccResponse)
		})
		It("GetContractIDs from Org included in Contract, expected to succeed", func() {
			ccResponse := (cc.From(userID1).Query(`GetContractIDs`))
			expectcc.ResponseOk(ccResponse)
		})
		It("GetContracts from Org included in Contract, expected to succeed", func() {
			ccResponse := (cc.From(userID1).Query(`GetContracts`))
			expectcc.ResponseOk(ccResponse)
		})
		It("UpdateContractItem from Org included in Contract, expected to succeed", func() {
			testID := "ID-01"
			ccResponse := (cc.From(userID1).Invoke(`UpdateContractItem`, testID, &payload.Item{
				Enabled:    &f,
				Write:      &f,
				ObjectId:   item1,
				OrgId:      org1,
				ObjectType: type1,
			}))
			expectcc.ResponseOk(ccResponse)
		})
		It("DeleteContractItem from Org included in Contract, expected to succeed", func() {
			testID := "ID-01"
			ccResponse := (cc.From(userID1).Invoke(`DeleteContractItem`, testID, item1))
			expectcc.ResponseOk(ccResponse)
		})
		It("AddContractItem from Org included in Contract, expected to succeed", func() {
			testID := "ID-04"
			ccResponse := (cc.From(userID4).Invoke(`AddContractItem`, testID, &payload.Item{
				Enabled:    &t,
				Write:      &f,
				ObjectId:   uuid.New().String(),
				OrgId:      org4,
				ObjectType: type4,
			}))
			expectcc.ResponseOk(ccResponse)
		})

	})

	Describe("Contract functions (checking handlers)", func() {

		It("ProposeContract with acl error, expected to fail", func() {
			ccResponse := (cc.From(maliciousUserID).Invoke(`ProposeContract`, &payload.ContractPayload{
				ContractId:   uuid.New().String(),
				ContractType: Contract_Type[rand.Intn(len(Contract_Type))],
				Orgs:         []string{org1, org2},
				Items: []payload.Item{{
					Enabled:    &t,
					Write:      &t,
					ObjectId:   uuid.New().String(),
					OrgId:      org1,
					ObjectType: type1,
				},
					{
						Enabled:    &t,
						Write:      &f,
						ObjectId:   uuid.New().String(),
						OrgId:      org2,
						ObjectType: type2,
					}},
			}))
			expectcc.ResponseError(ccResponse)
		})

		It("AcceptContract for rejected Contract, expected to fail", func() {

			testID := "ID-02"
			ccResponse := (cc.From(userID2).Invoke(`AcceptContract`, testID))
			expectcc.ResponseError(ccResponse)
		})

		It("RejectContract for accepted Contract, expected to fail", func() {

			testID := "ID-01"
			ccResponse := (cc.From(userID1).Invoke(`RejectContract`, testID))
			expectcc.ResponseError(ccResponse)
		})

		It("DeleteContract for non existing Contract, expected to fail", func() {
			testID := "ID-00000"
			ccResponse := (cc.From(userID3).Invoke(`DissolveContract`, testID))
			expectcc.ResponseError(ccResponse)
		})

		It("GetContractByID for Org with no Contracts, expected to fail", func() {
			testID := "ID-03"
			ccResponse := (cc.From(userID3).Query(`GetContractByID`, testID))
			expectcc.ResponseError(ccResponse)
		})

		It("GetContractIDs from Org with no Contracts. Returns empty array [], expected to succeed", func() {
			ccResponse := (cc.From(userID3).Query(`GetContractIDs`))
			expectcc.ResponseOk(ccResponse)
		})
		It("GetContracts from Org with no Contracts. Returns empty array [], expected to succeed", func() {
			ccResponse := (cc.From(userID1).Query(`GetContracts`))
			expectcc.ResponseOk(ccResponse)
		})
		It("UpdateContractItem for Item not included in the Contract, expected to fail", func() {
			testID := "ID-01"
			ccResponse := (cc.From(userID1).Invoke(`UpdateContractItem`, testID, &payload.Item{
				Enabled:    &f,
				Write:      &f,
				ObjectId:   item3,
				OrgId:      org1,
				ObjectType: type1,
			}))
			expectcc.ResponseError(ccResponse)
		})
		It("DeleteContractItem for Item not included in the Contract, expected to fail", func() {
			testID := "ID-01"
			ccResponse := (cc.From(userID1).Invoke(`DeleteContractItem`, testID, item3))
			expectcc.ResponseError(ccResponse)
		})
		It("AddContractItem from non accepted Contract, expected to fail", func() {
			testID := "ID-02"
			ccResponse := (cc.From(userID3).Invoke(`AddContractItem`, testID, &payload.Item{
				Enabled:    &t,
				Write:      &f,
				ObjectId:   uuid.New().String(),
				OrgId:      org4,
				ObjectType: type4,
			}))
			expectcc.ResponseError(ccResponse)
		})

	})
})
