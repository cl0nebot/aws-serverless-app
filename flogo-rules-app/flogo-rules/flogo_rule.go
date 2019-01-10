package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

// OrgStatus data
type OrgStatus struct {
	OrgID         string `json:"id"`
	Status        string `json:"status,omitempty"`
	EffectiveDate string `json:"effective,omitempty"`
	Inforce       bool   `json:"inforce,omitempty"`
}

func invokeRules(org *OrgStatus) {
	startTime := time.Now()
	log.Println("** Invoke flogo rules **")

	rs, err := startRuleSession("asession")
	if err != nil {
		log.Printf("Error starting rule session [%s]\n", err)
		panic(err)
	}

	// Now assert a "orgstatus" tuple
	log.Printf("Asserting orgstatus tuple of id %s\n", org.OrgID)
	t1, _ := model.NewTupleWithKeyValues("orgstatus", org.OrgID)
	t1.SetString(nil, "OrgID", org.OrgID)
	t1.SetString(nil, "Status", org.Status)
	t1.SetString(nil, "EffectiveDate", org.EffectiveDate)
	// default Inforce=true if no rule fires
	t1.SetBool(nil, "Inforce", true)
	rs.Assert(nil, t1)

	// get rule result and update the OrgStatus
	inforce, _ := t1.GetBool("Inforce")
	log.Printf("Set rule result inforce=%t\n", inforce)
	org.Inforce = inforce

	// Retract tuples
	rs.Retract(nil, t1)

	//unregister the session, i.e; cleanup
	rs.Unregister()
	log.Printf("Rule execution elapsed time %s\n", time.Since(startTime))
}

func startRuleSession(name string) (model.RuleSession, error) {
	log.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	// Register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		log.Printf("Error [%s]\n", err)
		return nil, err
	}

	// Create a RuleSession
	rs, _ := ruleapi.GetOrCreateRuleSession(name)

	// rule to check org's effective date
	rule := notEffectiveRule("Not Effective")
	rs.AddRule(rule)
	log.Printf("Rule added: [%s]\n", rule.GetName())

	// rule to check if org is active
	rule2 := notActiveRule("Not Active")
	rs.AddRule(rule2)
	log.Printf("Rule added: [%s]\n", rule2.GetName())

	// Start the rule session
	rs.Start(nil)
	return rs, nil
}

const tupleDescriptor = `
[
  {
    "name": "orgstatus",
    "properties": [
      {
        "name": "OrgID",
        "type": "string",
        "pk-index": 0
      },
      {
        "name": "Status",
        "type": "string"
      },
      {
        "name": "EffectiveDate",
        "type": "string"
      },
      {
        "name": "Inforce",
        "type": "bool"
      }
    ]
  }
]`

// rule to check if orgstatus has a later effective date
func notEffectiveRule(name string) model.MutableRule {
	rule := ruleapi.NewRule(name)
	rule.AddCondition("c1", []string{"orgstatus"}, notOrgEffective, nil)
	rule.SetAction(notEffectiveAction)
	return rule
}

// condition returns true if org's effective date (YYYY-mm-dd) is later than current date
func notOrgEffective(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	t1 := tuples["orgstatus"]
	if t1 == nil {
		log.Println("Should not get a nil tuple in notOrgEffective! This is an error")
		return false
	}
	effDate, _ := t1.GetString("EffectiveDate")
	log.Printf("notOrgEffective condition got %s\n", effDate)
	currDate := time.Now().Format("2006-01-02")
	return strings.Compare(currDate, effDate) < 0
}

// set org status not in-force because effective date is later than current date
func notEffectiveAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	log.Printf("Rule fired: [%s]\n", ruleName)
	t1 := tuples["orgstatus"].(model.MutableTuple)
	if t1 == nil {
		log.Println("Should not get nil tuples in notEffectiveAction! This is an error")
	} else {
		id, _ := t1.GetString("OrgID")
		log.Printf("Set org %s not in-force\n", id)
		t1.SetBool(ctx, "Inforce", false)
	}
}

// rule to check if orgstatus has status="active"
func notActiveRule(name string) model.MutableRule {
	rule := ruleapi.NewRule(name)
	rule.AddCondition("c1", []string{"orgstatus"}, notOrgActive, nil)
	rule.SetAction(notActiveAction)
	return rule
}

// condition returns true if org status is not "active" ignore case
func notOrgActive(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	t1 := tuples["orgstatus"]
	if t1 == nil {
		log.Println("Should not get nil tuples here in notOrgActive! This is an error")
		return false
	}
	id, _ := t1.GetString("OrgID")
	status, _ := t1.GetString("Status")
	log.Printf("notOrgActive condition got ID %s and status %s\n", id, status)
	return !strings.EqualFold("active", status)
}

// set org status not in-force because org status is not "active"
func notActiveAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	log.Printf("Rule fired: [%s]\n", ruleName)
	t1 := tuples["orgstatus"].(model.MutableTuple)
	if t1 == nil {
		log.Println("Should not get nil tuples here in notActiveAction! This is an error")
		return
	}
	id, _ := t1.GetString("OrgID")
	log.Printf("Set org %s not in-force\n", id)
	t1.SetBool(ctx, "Inforce", false)
}
