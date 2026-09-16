# Feature: Adapt CogOS Cognitive Loop (NREM/REM & Curiosity)
**Status:** `Approved`

## 1. Mục tiêu (Context & Goals)
Porting kiến trúc Graph-based cognitive từ `cogos` (Python) sang `cognitiveos` (Golang). Đặc biệt mô phỏng cơ chế Offline processing (NREM/REM) và Semantic Gap Detection để giúp hạt nhân tự nhận thức các lỗ hổng kiến thức (Curiosity Loop) dựa trên dữ liệu hiện tại trong `cogfield`.

## 2. Tiêu chí Nghiệm thu (Acceptance Criteria)
- [x] Định nghĩa `SemanticGap` và các constants `GapType` theo chuẩn CDLC.
- [x] Bổ sung function `DetectGaps` trong `pkg/substrate/cogfield/gaps.go` để tìm Structural Gaps, Missing Gaps và Semantic Gaps.
- [x] Có Unit tests đảm bảo cơ chế phát hiện Gap hoạt động chính xác.
- [x] Cập nhật tiến độ vào `.cdlc/ROADMAP.md`.

## 3. Đặc tả Kiến trúc & Data Model (Core A-SDD)
- **Structs/Interfaces/Types:**
```go
package cogfield

type GapType string

const (
	GapStructural GapType = "structural"
	GapSemantic   GapType = "semantic"
	GapMissing    GapType = "missing"
)

type SemanticGap struct {
	GapType           GapType
	SubjectID         string
	SubjectLabel      string
	EdgeID            string // Có thể rỗng nếu là Missing/Structural
	Relation          string
	ObjectLabel       string
	Severity          float64
	Explanation       string
	SuggestedQuestion string
}
```
- **Tương tác (System Impact):** Các module xử lý nội bộ (ví dụ: `autonomic_ticker`) có thể gọi `cogfield.DetectGaps(graph)` để phân tích tình trạng kiến thức và sinh ra báo cáo "hypotheses/questions" tương tự như REM phase trong `cogos`.

## 4. Kế hoạch Triển khai (Implementation Steps)
- **Step 1:** Tạo file `pkg/substrate/cogfield/gaps.go` chứa logic `DetectGaps`.
- **Step 2:** Khớp với heuristic của Python:
  - *Structural:* Node strength thấp (ví dụ: < 0.45) nhưng có kết nối đến (incoming edges) và không có kết nối đi (outgoing edges).
  - *Semantic:* Edge có Weight cao nhưng thiếu path trung gian để giải thích.
  - *Missing:* (Dễ bị nhiễu) - có thể scan Tags nếu không tìm thấy NodeID tương ứng.
- **Step 3:** Tạo file test `pkg/substrate/cogfield/gaps_test.go`.

## 5. Rủi ro & Edge Cases (Error Handling)
- Đồ thị lớn sẽ gây nghẽn O(V + E) trên main thread. Cần lưu ý không chạy `DetectGaps` ở những nơi blocking I/O, mà chỉ nên trigger trong background (như Autonomous Scheduler - NREM/REM cycle).
