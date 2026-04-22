package mongo_model

import (
	"app/domain/model"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDefaultFilter_Query_WithID(t *testing.T) {
	id := primitive.NewObjectID()
	f := DefaultFilter{ID: id}
	query := make(map[string]any)
	f.Query(query)

	if query["_id"] != id {
		t.Errorf("_id = %v, want %v", query["_id"], id)
	}
}

func TestDefaultFilter_Query_WithIDStr(t *testing.T) {
	id := primitive.NewObjectID()
	idStr := id.Hex()
	f := DefaultFilter{IDStr: &idStr}
	query := make(map[string]any)
	f.Query(query)

	if query["_id"] != id {
		t.Errorf("_id = %v, want %v", query["_id"], id)
	}
}

func TestDefaultFilter_Query_WithIDs(t *testing.T) {
	ids := []primitive.ObjectID{primitive.NewObjectID(), primitive.NewObjectID()}
	f := DefaultFilter{IDs: ids}
	query := make(map[string]any)
	f.Query(query)

	idQuery, ok := query["_id"].(map[string]any)
	if !ok {
		t.Fatalf("_id query type = %T, want map[string]any", query["_id"])
	}
	inIDs, ok := idQuery["$in"].([]primitive.ObjectID)
	if !ok {
		t.Fatalf("$in type = %T, want []primitive.ObjectID", idQuery["$in"])
	}
	if len(inIDs) != 2 {
		t.Errorf("$in length = %d, want 2", len(inIDs))
	}
}

func TestDefaultFilter_Query_WithIDsStr(t *testing.T) {
	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	f := DefaultFilter{IDsStr: []string{id1.Hex(), id2.Hex()}}
	query := make(map[string]any)
	f.Query(query)

	idResult, ok := query["_id"].([]primitive.ObjectID)
	if !ok {
		t.Fatalf("_id type = %T, want []primitive.ObjectID", query["_id"])
	}
	if len(idResult) != 2 {
		t.Errorf("_id length = %d, want 2", len(idResult))
	}
}

func TestDefaultFilter_Query_WithIDsStr_InvalidHex(t *testing.T) {
	f := DefaultFilter{IDsStr: []string{"invalid-hex", primitive.NewObjectID().Hex()}}
	query := make(map[string]any)
	f.Query(query)

	idResult, ok := query["_id"].([]primitive.ObjectID)
	if !ok {
		t.Fatalf("_id type = %T, want []primitive.ObjectID", query["_id"])
	}
	// Only the valid hex should be included
	if len(idResult) != 1 {
		t.Errorf("_id length = %d, want 1 (only valid IDs)", len(idResult))
	}
}

func TestDefaultFilter_Query_CreatedAtGt(t *testing.T) {
	now := time.Now()
	f := DefaultFilter{CreatedAtGt: &now}
	query := make(map[string]any)
	f.Query(query)

	createdAt, ok := query["createdAt"].(map[string]any)
	if !ok {
		t.Fatalf("createdAt type = %T, want map[string]any", query["createdAt"])
	}
	if createdAt["$gt"] != &now {
		t.Error("$gt should be set")
	}
}

func TestDefaultFilter_Query_CreatedAtRange(t *testing.T) {
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	f := DefaultFilter{
		CreatedAtRange: &model.DatetimeRange{Start: start, End: end},
	}
	query := make(map[string]any)
	f.Query(query)

	createdAt, ok := query["createdAt"].(map[string]any)
	if !ok {
		t.Fatalf("createdAt type = %T, want map[string]any", query["createdAt"])
	}
	if createdAt["$gte"] != start {
		t.Error("$gte should equal start time")
	}
	if createdAt["$lte"] != end {
		t.Error("$lte should equal end time")
	}
}

func TestDefaultFilter_Query_UpdatedAtGt(t *testing.T) {
	now := time.Now()
	f := DefaultFilter{UpdatedAtGt: &now}
	query := make(map[string]any)
	f.Query(query)

	updatedAt, ok := query["updatedAt"].(map[string]any)
	if !ok {
		t.Fatalf("updatedAt type = %T, want map[string]any", query["updatedAt"])
	}
	if updatedAt["$gt"] != &now {
		t.Error("$gt should be set")
	}
}

func TestDefaultFilter_Query_UpdatedAtRange(t *testing.T) {
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()
	f := DefaultFilter{
		UpdatedAtRange: &model.DatetimeRange{Start: start, End: end},
	}
	query := make(map[string]any)
	f.Query(query)

	updatedAt, ok := query["updatedAt"].(map[string]any)
	if !ok {
		t.Fatalf("updatedAt type = %T, want map[string]any", query["updatedAt"])
	}
	if updatedAt["$gte"] != start {
		t.Error("$gte should equal start time")
	}
	if updatedAt["$lte"] != end {
		t.Error("$lte should equal end time")
	}
}

func TestDefaultFilter_Query_Raw(t *testing.T) {
	f := DefaultFilter{
		Raw: map[string]any{
			"status":  "active",
			"visible": true,
		},
	}
	query := make(map[string]any)
	f.Query(query)

	if query["status"] != "active" {
		t.Errorf("status = %v, want %q", query["status"], "active")
	}
	if query["visible"] != true {
		t.Errorf("visible = %v, want true", query["visible"])
	}
}

func TestDefaultFilter_Query_Empty(t *testing.T) {
	f := DefaultFilter{}
	query := make(map[string]any)
	f.Query(query)

	if len(query) != 0 {
		t.Errorf("empty filter should produce empty query, got %v", query)
	}
}

func TestDefaultFilter_FindOptions_WithLimit(t *testing.T) {
	limit := int64(10)
	f := DefaultFilter{Limit: &limit}
	opts := f.FindOptions()

	if opts.Limit == nil || *opts.Limit != 10 {
		t.Errorf("Limit = %v, want 10", opts.Limit)
	}
}

func TestDefaultFilter_FindOptions_WithOffset(t *testing.T) {
	offset := int64(20)
	f := DefaultFilter{Offset: &offset}
	opts := f.FindOptions()

	if opts.Skip == nil || *opts.Skip != 20 {
		t.Errorf("Skip = %v, want 20", opts.Skip)
	}
}

func TestDefaultFilter_FindOptions_WithSorts(t *testing.T) {
	f := DefaultFilter{
		Sorts: bson.D{{Key: "created_at", Value: -1}},
	}
	opts := f.FindOptions()

	if opts.Sort == nil {
		t.Error("Sort should not be nil")
	}
}

func TestDefaultFilter_FindOptions_Empty(t *testing.T) {
	f := DefaultFilter{}
	opts := f.FindOptions()

	if opts.Limit != nil {
		t.Error("Limit should be nil")
	}
	if opts.Skip != nil {
		t.Error("Skip should be nil")
	}
	if opts.Sort != nil {
		t.Error("Sort should be nil")
	}
}

func TestDefaultFilter_FindOptions_AllOptions(t *testing.T) {
	limit := int64(5)
	offset := int64(10)
	f := DefaultFilter{
		Limit:  &limit,
		Offset: &offset,
		Sorts:  bson.D{{Key: "name", Value: 1}},
	}
	opts := f.FindOptions()

	if opts.Limit == nil || *opts.Limit != 5 {
		t.Errorf("Limit = %v, want 5", opts.Limit)
	}
	if opts.Skip == nil || *opts.Skip != 10 {
		t.Errorf("Skip = %v, want 10", opts.Skip)
	}
	if opts.Sort == nil {
		t.Error("Sort should not be nil")
	}
}
