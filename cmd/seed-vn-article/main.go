package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
)

func main() {
	creds := flag.String("creds", "golang-blogs-firebase-adminsdk-fbsvc-64dacce61f.json", "Firebase service account JSON")
	project := flag.String("project", "golang-blogs", "Firebase project ID")
	collection := flag.String("collection", "articles", "Firestore collection name")
	flag.Parse()

	visible := true
	article := map[string]any{
		"id":          "timeline-thong-bao-tai-san-dien-tu-vn-2026",
		"slug":        "timeline-thong-bao-tai-san-dien-tu-vn-2026",
		"title":       "Timeline 2026: Các thông báo của Chính phủ Việt Nam về thị trường tài sản điện tử",
		"summary":     "Tổng hợp các mốc chính từ đầu năm 2026 đến tháng 6: hồ sơ cấp phép sàn, Quyết định 96, tiến độ thí điểm và kỳ vọng vận hành thị trường trong quý III.",
		"bodyHtml":    bodyHTML(),
		"category":    "Việt Nam",
		"author":      "Tony Blogs",
		"publishedAt": "2026-06-12T08:00:00Z",
		"tags":        []string{"vietnam", "tai-san-dien-tu", "crypto", "nq-05", "btc"},
		"thumbnail":   "https://images.unsplash.com/photo-1611974789855-9c2a0a7236a3?auto=format&fit=crop&w=1200&q=80",
		"visible":     visible,
	}

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, *project, option.WithCredentialsFile(*creds))
	if err != nil {
		log.Fatalf("firestore client: %v", err)
	}
	defer client.Close()

	id := article["id"].(string)
	if _, err := client.Collection(*collection).Doc(id).Set(ctx, article); err != nil {
		log.Fatalf("write: %v", err)
	}

	fmt.Printf("uploaded %s -> %s/%s\n", id, *project, *collection)
}

func bodyHTML() string {
	return `<p>Sau khi <strong>Nghị quyết 05/2025/NQ-CP</strong> có hiệu lực từ tháng 9/2025, năm 2026 là giai đoạn triển khai thực chất của chương trình <em>thí điểm 5 năm</em> thị trường tài sản mã hóa tại Việt Nam. Dưới đây là các mốc thông báo và chính sách nổi bật từ đầu năm đến tháng 6/2026.</p>

<h2>Timeline chính (01/2026 – 06/2026)</h2>

<ul>
<li><strong>09/01/2026</strong> — Bộ Tài chính và Ủy ban Chứng khoán Nhà nước (UBCKNN) công bố tiếp nhận hồ sơ đăng ký cấp phép dịch vụ tổ chức thị trường giao dịch tài sản mã hóa theo cơ chế nộp hồ sơ theo từng giai đoạn.</li>
<li><strong>15/01/2026</strong> — Hạn chót nộp hồ sơ đăng ký cấp phép sàn giao dịch trong đợt đầu. Theo các báo cáo ngành, có <strong>7 tổ chức</strong> đã nộp hồ sơ sơ bộ.</li>
<li><strong>20/01/2026</strong> — <strong>Quyết định 96/QĐ-BTC</strong> ban hành, quy định chi tiết thủ tục, hồ sơ, thời hạn xử lý cấp phép, sửa đổi và thu hồi giấy phép cung cấp dịch vụ tổ chức thị trường giao dịch tài sản mã hóa.</li>
<li><strong>28/02/2026</strong> — Mốc kỳ vọng khởi động giai đoạn vận hành thí điểm ban đầu; tuy nhiên đến cuối quý I, <strong>chưa có sàn nào được cấp phép chính thức</strong> do yêu cầu vốn, an ninh mạng và thẩm định liên ngành.</li>
<li><strong>03/2026 – 04/2026</strong> — UBCKNN tiếp tục rà soát hồ sơ bổ sung. Một số ứng viên rút lui hoặc đang huy động vốn để đáp ứng ngưỡng vốn điều lệ theo khung thí điểm. Thị trường chờ tín hiệu cấp phép đầu tiên.</li>
<li><strong>05/2026</strong> — Bộ Tài chính nhấn mạnh nguyên tắc <strong>giao dịch qua nhà cung cấp được cấp phép</strong>, thanh toán bằng VND, tuân thủ phòng chống rửa tiền và an ninh thông tin cấp độ 4. Cảnh báo giao dịch qua kênh không được phép vẫn có rủi ro pháp lý.</li>
<li><strong>06/2026 (hiện tại)</strong> — Thí điểm tiếp tục; kỳ vọng phần lớn nhà đầu tư theo dõi <strong>quý III/2026</strong> là cửa sổ có thể ghi nhận sàn đầu tiên vận hành, sau đó nhà đầu tư trong nước có thời gian chuyển đổi sang nền tảng hợp pháp theo lộ trình chuyển tiếp.</li>
</ul>

<h2>Điểm nhà đầu tư cần nhớ</h2>

<ol>
<li>Chỉ tối đa <strong>5 giấy phép</strong> sàn trong giai đoạn thí điểm theo Nghị quyết 05.</li>
<li>Hồ sơ cấp phép yêu cầu thẩm định của Bộ Công An về an toàn hệ thống thông tin, phối hợp NHNN về phòng chống rửa tiền.</li>
<li>Giao dịch, thanh toán, phát hành liên quan tài sản mã hóa trong khung thí điểm phải tuân thủ quy định về <strong>đồng Việt Nam</strong>.</li>
<li>Sau khi có sàn được cấp phép, nhà đầu tư trong nước sẽ có <strong>thời gian chuyển tiếp</strong> để chuyển hoạt động sang nền tảng được phép.</li>
</ol>

<blockquote><p>Đây là bản tổng hợp timeline phục vụ đọc nhanh, không thay thế văn bản pháp luật. Theo dõi cổng Bộ Tài chính và UBCKNN để cập nhật thông báo chính thức.</p></blockquote>

<p><em>Nguồn tham khảo công khai: Nghị quyết 05/2025/NQ-CP, Quyết định 96/QĐ-BTC và thông tin tiếp nhận hồ sơ của UBCKNN.</em></p>`
}
