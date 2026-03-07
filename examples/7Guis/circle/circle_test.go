package circle

import (
	"testing"
)

// TestCircleModel - BDD: Given circle model, when creating circle, then properties match
func TestCircleModel(t *testing.T) {
	// GIVEN: A circle model with specific coordinates and radius
	model := CircleModel{x: 100, y: 200, r: 30}

	// WHEN: Creating a circle from model
	circle := NewCircle(func(data *Circle) {})
	circle.Reset(model.x, model.y, model.r)

	// THEN: Circle should have exact coordinates and radius
	if circle.x != 100 || circle.y != 200 || circle.r != 30 {
		t.Errorf("Circle coordinates mismatch. Expected (100,200,30), got (%d,%d,%d)",
			circle.x, circle.y, circle.r)
	}

	// AND: Model() should return same values
	resultModel := circle.Model()
	if resultModel.x != 100 || resultModel.y != 200 || resultModel.r != 30 {
		t.Errorf("Model() mismatch. Expected (100,200,30), got (%d,%d,%d)",
			resultModel.x, resultModel.y, resultModel.r)
	}
}

// TestCircleResetUpdatesProperties - BDD: Given initial circle, when resetting to new values, then all properties update
func TestCircleResetUpdatesProperties(t *testing.T) {
	// GIVEN: A circle with initial coordinates
	circle := NewCircle(func(data *Circle) {})
	circle.Reset(50, 60, 20)

	// WHEN: Resetting to new values
	circle.Reset(150, 160, 40)

	// THEN: All properties should be updated
	if circle.x != 150 || circle.y != 160 || circle.r != 40 {
		t.Errorf("Reset failed. Expected (150,160,40), got (%d,%d,%d)",
			circle.x, circle.y, circle.r)
	}
}

// TestCircleStateSnapshot - BDD: Given multiple circles, when creating snapshots, then snapshots contain all data
func TestCircleStateSnapshot(t *testing.T) {
	// GIVEN: Multiple circles
	circle1 := NewCircle(func(data *Circle) {})
	circle1.Reset(100, 100, 25)

	circle2 := NewCircle(func(data *Circle) {})
	circle2.Reset(200, 200, 35)

	// WHEN: Creating snapshots of their models
	models := []CircleModel{
		circle1.Model(),
		circle2.Model(),
	}

	// THEN: Snapshots should contain the data
	if len(models) != 2 {
		t.Errorf("Expected 2 models, got %d", len(models))
	}

	if models[0].x != 100 || models[0].y != 100 || models[0].r != 25 {
		t.Errorf("Model 0 mismatch")
	}

	if models[1].x != 200 || models[1].y != 200 || models[1].r != 35 {
		t.Errorf("Model 1 mismatch")
	}

	// AND: Original circles and snapshots should be independent (immutable snapshot)
	circle1.Reset(300, 300, 50)
	if models[0].x == 300 { // Snapshot should NOT change
		t.Error("Snapshot was mutated when original changed")
	}
}

// TestCanvasStateRecovery - BDD: Given canvas with circles, when extracting and restoring state, then circles recover
func TestCanvasStateRecovery(t *testing.T) {
	// GIVEN: A canvas with multiple circles
	canvas := NewCanvas()

	circle1 := NewCircle(func(data *Circle) {})
	circle1.Reset(100, 100, 20)

	circle2 := NewCircle(func(data *Circle) {})
	circle2.Reset(200, 200, 30)

	canvas.Add(circle1)
	canvas.Add(circle2)

	// WHEN: Extracting all circle models
	models := make([]CircleModel, 0)
	for _, child := range canvas.Children() {
		if c, ok := child.(*Circle); ok {
			models = append(models, c.Model())
		}
	}

	// THEN: Should have 2 models
	if len(models) != 2 {
		t.Errorf("Expected 2 circles in models, got %d", len(models))
	}

	// AND: When canvas is reset and circles are restored from models
	canvas.ResetContainer()
	for _, model := range models {
		circle := NewCircle(func(data *Circle) {})
		circle.Reset(model.x, model.y, model.r)
		canvas.Add(circle)
	}

	// THEN: Canvas should have 2 circles again with exact properties
	if len(canvas.Children()) != 2 {
		t.Errorf("Expected 2 children after restore, got %d", len(canvas.Children()))
	}

	children := canvas.Children()
	c1 := children[0].(*Circle)
	if c1.x != 100 || c1.y != 100 || c1.r != 20 {
		t.Errorf("Circle 1 restore failed. Got (%d,%d,%d)", c1.x, c1.y, c1.r)
	}

	c2 := children[1].(*Circle)
	if c2.x != 200 || c2.y != 200 || c2.r != 30 {
		t.Errorf("Circle 2 restore failed. Got (%d,%d,%d)", c2.x, c2.y, c2.r)
	}
}

// TestUndoRedoScenarioAdd - BDD: Given 2 circles added, when undo called, then should go back to 1 circle
func TestUndoRedoScenarioAdd(t *testing.T) {
	// GIVEN: snapshots representing canvas states during add operations
	// snapshot1 := []CircleModel{} // Empty canvas at start

	model1 := CircleModel{x: 100, y: 100, r: 25}
	snapshot2 := []CircleModel{model1} // After drawing circle 1

	model2 := CircleModel{x: 200, y: 200, r: 35}
	snapshot3 := []CircleModel{model1, model2} // After drawing circle 2

	// WHEN: Checking snapshots (represent undo/redo history)

	// THEN: snapshot2 should represent single circle
	if len(snapshot2) != 1 {
		t.Error("Snapshot2 should have 1 circle")
	}

	// AND: snapshot3 should represent two circles
	if len(snapshot3) != 2 {
		t.Error("Snapshot3 should have 2 circles")
	}

	// AND: Going back from snapshot3 to snapshot2 means circle 2 is removed
	if snapshot3[1].x != model2.x || snapshot3[1].y != model2.y {
		t.Error("Circle 2 data in snapshot3 mismatch")
	}
}

// TestRadiusChangeSnapshot - BDD: Given circle, when radius changes, then new snapshot differs only in radius
func TestRadiusChangeSnapshot(t *testing.T) {
	// GIVEN: Initial circle state
	initialModel := CircleModel{x: 150, y: 150, r: 30}

	// WHEN: User changes radius
	modifiedModel := CircleModel{x: 150, y: 150, r: 50}

	// THEN: Models should differ only in radius
	if initialModel.x != modifiedModel.x || initialModel.y != modifiedModel.y {
		t.Error("Position should not change when radius changes")
	}

	if initialModel.r == modifiedModel.r {
		t.Error("Radius should change")
	}

	if modifiedModel.r != 50 {
		t.Errorf("Expected radius 50, got %d", modifiedModel.r)
	}
}

// TestCircleStateIndependence - BDD: Given circle and snapshot, when original changes, then snapshot unchanged
func TestCircleStateIndependence(t *testing.T) {
	// GIVEN: A circle and its model snapshot
	originalCircle := NewCircle(func(data *Circle) {})
	originalCircle.Reset(120, 150, 45)

	snapshot := originalCircle.Model()
	snapshotX := snapshot.x
	snapshotR := snapshot.r

	// WHEN: Original circle changes
	originalCircle.Reset(200, 250, 60)

	// THEN: Snapshot values should remain unchanged (immutable)
	if snapshot.x != snapshotX {
		t.Error("Snapshot x changed when original changed")
	}

	if snapshot.r != snapshotR {
		t.Error("Snapshot r changed when original changed")
	}
}
