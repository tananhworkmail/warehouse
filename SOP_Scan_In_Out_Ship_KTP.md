# SOP vận hành Scan In, Scan Out và Scan Ship kho thành phẩm KTP

| Thông tin | Nội dung |
|---|---|
| Mã tài liệu | SOP-KTP-SCAN-001 |
| Phiên bản | 1.0 |
| Ngày ban hành | 25/08/2026 |
| Phạm vi | Kho thành phẩm KTP |
| Hệ thống | Warehouse Web và ERP LYS |

## 1. Mục đích

Quy định cách vận hành ba chức năng:

- **Scan In:** nhập thùng thành phẩm vào kho KTP hoặc tái nhập sau Tái chế/Kiểm hàng.
- **Scan Out:** đưa thùng đang ở KTP sang luồng **Tái chế** hoặc **Kiểm hàng**.
- **Scan Ship:** xuất thùng từ KTP lên xe/container để giao khách hàng.

Mục tiêu là mỗi thùng được kiểm tra đúng điều kiện, cập nhật ngay khi quét và có thể truy vết người thao tác, thời gian, vị trí và trọng lượng.

## 2. Trách nhiệm

### 2.1. Nhân viên vận hành

- Đăng nhập bằng tài khoản của chính mình.
- Kiểm tra đúng thùng, chuyền, vị trí kho, container và biển số trước khi quét.
- Chỉ di chuyển thùng sau khi màn hình báo thành công.
- Không tự xác nhận lệch cân nếu chưa kiểm tra lại thùng và cân.
- Báo quản lý khi hệ thống không tìm thấy carton, sai trạng thái, sai Shipping Plan hoặc lỗi database.

### 2.2. Quản lý kho

- Kiểm soát quyền sử dụng từng chức năng.
- Kiểm tra và xử lý các trường hợp xác nhận lệch cân.
- Đối chiếu số thùng, số đôi và trạng thái ERP khi có chênh lệch.

### 2.3. IT/ERP

- Đảm bảo kết nối Warehouse Web, API, ERP, LYS_WEB, máy quét và cân điện tử.
- Xử lý lỗi tài khoản, quyền, database và dữ liệu nguồn.
- Hỗ trợ đối chiếu log xác nhận lệch cân.

## 3. Thiết bị và điều kiện trước khi làm việc

- Máy tính có kết nối mạng nội bộ.
- Máy quét mã vạch hoạt động bình thường.
- Cân điện tử bắt buộc đối với **Scan In** và **Scan Ship**.
- Tem `CARTONBAR` hoặc `LPNBarCode` đọc được.
- Tài khoản có quyền đúng với màn hình Scan In, Scan Out hoặc Scan Ship.
- Không dùng chung tài khoản giữa nhiều người.

Hệ thống nhận mã theo thứ tự:

1. Tìm đúng `YWCP.CARTONBAR`.
2. Nếu không thấy, tìm `YWCP.LPNBarCode`; với chuỗi dài, hệ thống dùng 15 ký tự cuối để đối chiếu LPN.

## 4. Quy ước trạng thái carton

| `YWCP.SB` | Ý nghĩa trong quy trình KTP | Chức năng tiếp theo được phép |
|---|---|---|
| `NULL` hoặc rỗng | Chưa nhập kho thành phẩm | Scan In |
| `1` | Đang ở kho thành phẩm | Scan Out hoặc Scan Ship |
| `2` | Đã đưa đi Tái chế | Scan In để tái nhập |
| `4` | Đã đưa đi Kiểm hàng | Scan In để tái nhập |
| `3` | Đã xuất giao hàng | Không được Scan In/Out/Ship lại |

Luồng trạng thái chuẩn:

```text
Carton mới (SB rỗng) ──Scan In──> SB=1 ──Scan Ship────────> SB=3
                                      │
                                      ├──Scan Out: Tái chế──> SB=2 ──Scan In──> SB=1
                                      │
                                      └──Scan Out: Kiểm hàng─> SB=4 ──Scan In──> SB=1
```

## 5. Quy tắc chung cho cả ba chức năng

1. Mỗi carton quét thành công được cập nhật database **ngay lập tức**; không có bước Upload hoặc cập nhật hàng loạt.
2. Danh sách phiên chỉ dùng để xem các carton vừa thao tác, số thùng và số đôi; danh sách này không quyết định carton có được quét hay không.
3. Được phép quét carton khác RY trong cùng phiên. Hệ thống sẽ hiện cảnh báo **khác đơn hàng/khác RY** nhưng vẫn tiếp tục xử lý nếu carton đạt các điều kiện khác.
4. Carton đã có trong danh sách của phiên hiện tại không được quét trùng trên giao diện.
5. Nút **Refresh phiên/Kết thúc phiên/Phiên mới** chỉ quản lý danh sách trên màn hình; không hoàn tác các carton đã ghi vào ERP.
6. Không quét lại một carton chỉ vì đã Refresh danh sách. Trạng thái thật nằm trong ERP và vẫn được backend kiểm tra lại.
7. Chỉ công nhận thành công khi có thông báo thành công và carton xuất hiện trong danh sách phiên.
8. Nếu phiên đăng nhập hết hạn (`401`), đăng nhập lại ở đúng trang đăng nhập của chức năng đang dùng.

## 6. SOP Scan In

### 6.1. Mục đích

Nhập carton thành phẩm vào kho KTP lần đầu hoặc tái nhập carton từ luồng Tái chế/Kiểm hàng.

### 6.2. Chuẩn bị màn hình

1. Mở trang **Scan In KTP** và đăng nhập.
2. Chọn chế độ:
   - **Carton:** chỉ quét mã thùng.
   - **Carton and Inbox:** quét mã thùng, sau đó quét đủ các hộp con của thùng.
3. Chọn **Production Line**. Đây là thông tin bắt buộc và phải có trong danh mục chuyền hợp lệ của ERP.
4. Nhập/chọn **Storage**. Đây là thông tin bắt buộc; giá trị mặc định trên giao diện là `A2`.
5. Nhập **Location** nếu có. Trường này được phép để trống.
6. Kiểm tra cân đã tự kết nối và trọng lượng đang hiển thị lớn hơn `0`.
7. Đặt đúng một carton lên cân, chờ số cân ổn định rồi mới quét.

Kho được hệ thống cố định là `A`; người vận hành không phải chọn kho.

### 6.3. Các bước quét chế độ Carton

1. Đặt carton lên cân.
2. Kiểm tra số cân đã ổn định.
3. Quét `CARTONBAR` hoặc LPN.
4. Hệ thống đọc carton và kiểm tra điều kiện ở mục 6.5.
5. Nếu carton khác RY với các carton trong danh sách phiên, đọc cảnh báo và kiểm tra lại tem; nếu đúng carton cần nhập thì không phải dừng phiên.
6. Nếu tất cả điều kiện đạt, hệ thống tự ghi Scan In ngay, không cần bấm Upload.
7. Xác nhận màn hình báo thành công; kiểm tra carton đã vào danh sách và số thùng/số đôi tăng đúng.
8. Nhấc carton khỏi cân, chờ cân về trạng thái sẵn sàng rồi thực hiện carton tiếp theo.

### 6.4. Các bước quét chế độ Carton and Inbox

1. Thực hiện các bước chuẩn bị và quét mã carton như chế độ Carton.
2. Khi carton đạt các điều kiện chính, màn hình chuyển sang nhận mã hộp con.
3. Quét lần lượt từng mã Inbox đúng số lượng yêu cầu.
4. Không được quét mã không thuộc carton hoặc quét một mã quá số lượng quy định.
5. Khi đã quét đủ toàn bộ Inbox, hệ thống tự ghi Scan In cho carton.
6. Xác nhận thông báo thành công và dòng carton trong danh sách phiên.

Nếu thiếu dữ liệu Inbox, thiếu mã, thừa mã hoặc sai mã thì carton chưa được ghi vào kho.

### 6.5. Điều kiện Scan In

Carton chỉ được Scan In khi đồng thời đạt các nhóm điều kiện sau:

#### A. Carton và trạng thái

- Tìm thấy carton trong `YWCP`.
- Carton có thông tin đơn hàng/RY.
- `SB` đang rỗng, `2` hoặc `4`.
- `SB=1` bị từ chối vì carton đang ở KTP.
- `SB=3` bị từ chối vì carton đã xuất giao hàng.

#### B. Cân nặng

- Trọng lượng tiêu chuẩn thùng `sgw > 0`.
- Trọng lượng thực tế từ cân phải lớn hơn `0`.
- Hệ thống phải đọc được `Gram > 0`.
- Công thức kiểm tra:

```text
Sai số cho phép = Gram / 3
Đạt cân khi ABS(sgw - trọng lượng thực tế) <= Sai số cho phép
```

Nguồn `Gram` hiện tại:

- Nhà máy `VB1` hoặc `VB7`: `ROUND(sgw × 0,03; 3)`.
- Nhà máy khác: `YWXXNew.NWeigh × 1,04` theo model/màu/size.

#### C. Sản lượng Gò

- Chỉ khóa sản lượng đối với đơn thuộc nhà máy có phân loại `PD`.
- Sản lượng sau khi nhập carton không được lớn hơn sản lượng Gò hoàn thành của đúng đơn/size.
- Nếu ERP có cờ cho phép `Flag=1`, hệ thống cho qua kiểm tra sản lượng.

#### D. Cấp phát vật tư

- Phải có dữ liệu đơn vật tư, số đôi đơn hàng và danh sách vật tư bắt buộc.
- Không được có vật tư bắt buộc hoàn toàn chưa mở phiếu cấp phát.
- Với carton nhập lần đầu, mức vật tư đã cấp phải bao phủ được số lượng sau carton hoặc ERP có `Flag=1`.
- Với carton tái nhập từ `SB=2/4`, hệ thống không khóa lại theo mức bao phủ sau khi đã xác nhận không có dòng vật tư bắt buộc bị thiếu hoàn toàn.
- Nếu dữ liệu vật tư cao su cần dùng nhưng không đọc được nguồn `LIY_DD`, carton bị từ chối.

#### E. Inbox

- Chỉ áp dụng với chế độ **Carton and Inbox**.
- Mỗi size có số lượng phải có mã Inbox hợp lệ.
- Mã quét phải thuộc đúng carton.
- Số lần quét từng mã không được vượt yêu cầu.
- Phải quét đủ toàn bộ số Inbox dự kiến.

### 6.6. Xử lý lệch cân Scan In

Chỉ được xác nhận lệch cân khi carton đạt điều kiện trạng thái, sản lượng, vật tư và Inbox; riêng cân nằm ngoài giới hạn.

1. Kiểm tra lại tem carton, vị trí đặt hàng và số cân.
2. Nếu không tiếp tục, chọn **Không/Hủy**, lấy carton ra kiểm tra; hệ thống không ghi Scan In.
3. Nếu được phép tiếp tục, chọn **Có/Xác nhận**.
4. Nhập đúng userid của tài khoản đang đăng nhập vào ô xác nhận.
5. Hệ thống kiểm tra userid với JWT, ghi carton và ghi log xác nhận lệch cân vào `LYS_WEB.dbo.KTP_Scan_logs`.

Không được dùng xác nhận lệch cân để bỏ qua lỗi trạng thái, sản lượng, vật tư hoặc Inbox.

### 6.7. Database được cập nhật khi Scan In thành công

| Bảng/field | Giá trị cập nhật |
|---|---|
| `YWCP.SB` | `1` |
| `YWCP.CKBH` | `A` |
| `YWCP.DepNo` | Production Line ở lần nhập đầu; tái nhập giữ chuyền cũ |
| `YWCP.KCBH` | Storage đã chọn |
| `YWCP.KVBH` | Location đã nhập, có thể rỗng |
| `YWCP.INCS` | Lần đầu là `1`; các lần sau tăng `1` |
| `YWCP.INDATE` | Ghi thời gian ở lần nhập đầu; tái nhập giữ ngày đầu |
| `YWCP.LastInDate` | Thời gian Scan In gần nhất |
| `YWCP.rgw` | Trọng lượng thực tế khi nhập |
| `YWCP.INUser` | Userid đăng nhập |
| `YWCP.REDATE` | Thời gian tái nhập nếu trạng thái trước là `2/4` |
| `YWCP.RECS` | Tăng `1` nếu trạng thái trước là `2/4` |
| `LYS_WEB.dbo.KTP_Scan_logs` | Chỉ thêm log khi có xác nhận lệch cân |

### 6.8. Kết thúc phiên Scan In

1. Đối chiếu số thùng và số đôi trên màn hình với số thực tế.
2. Bấm **Kết thúc phiên** khi hoàn tất.
3. Danh sách phiên, Production Line và Location được xóa để bắt đầu phiên mới; thao tác này không xóa dữ liệu đã Scan In trong ERP.

## 7. SOP Scan Out

### 7.1. Mục đích

Đưa carton đang ở kho thành phẩm sang một trong hai luồng:

- **Tái chế:** cập nhật `SB=2`.
- **Kiểm hàng:** cập nhật `SB=4`.

### 7.2. Các bước thao tác

1. Mở trang **Scan Out KTP** và đăng nhập.
2. Chọn chính xác một mục **Tái chế** hoặc **Kiểm hàng** trước khi quét.
3. Quét `CARTONBAR` hoặc LPN.
4. Hệ thống kiểm tra carton và trạng thái trong ERP.
5. Nếu carton khác RY với danh sách phiên, hệ thống cảnh báo nhưng vẫn cho xử lý nếu carton hợp lệ.
6. Khi đạt điều kiện, hệ thống cập nhật ngay carton theo lựa chọn, không cần bấm Upload.
7. Kiểm tra thông báo thành công, dòng carton mới và bộ đếm số thùng/số đôi.
8. Tiếp tục carton kế tiếp.

Lựa chọn Tái chế/Kiểm hàng được áp dụng tại thời điểm quét. Trước khi đổi mục, phải kiểm tra carton tiếp theo thuộc đúng luồng.

### 7.3. Điều kiện Scan Out

- Tìm thấy carton trong `YWCP`.
- Carton có RY/đơn hàng.
- Bắt buộc `SB=1`.
- `EXEDATE` phải rỗng, nghĩa là carton chưa Scan Ship.
- `SB` rỗng, `2`, `4`, `3` hoặc giá trị khác đều bị từ chối.

### 7.4. Database được cập nhật khi Scan Out thành công

| Lựa chọn | Cập nhật chính |
|---|---|
| Tái chế | `YWCP.SB = '2'` |
| Kiểm hàng | `YWCP.SB = '4'` |

Các cập nhật chung:

| Bảng/field | Giá trị cập nhật |
|---|---|
| `YWCP.OUTCS` | Tăng `1` |
| `YWCP.OUTDATE` | Thời gian Scan Out hiện tại |
| `DDZL.YN` | Chuyển từ `5` sang `2` đối với đơn liên quan |
| `YWDD.SB` | Chuyển từ `OK` sang `IN` đối với đơn đóng gói liên quan |

### 7.5. Kết thúc/Refresh phiên Scan Out

1. Đối chiếu số thùng và số đôi.
2. Có thể chọn từng carton trong danh sách để xem chi tiết đơn/size.
3. Bấm **Refresh phiên** để xóa danh sách vừa quét trên màn hình.
4. Refresh không hoàn tác `SB`, `OUTCS`, `OUTDATE` hoặc các cập nhật đơn hàng đã ghi.

## 8. SOP Scan Ship

### 8.1. Mục đích

Xuất carton đang ở kho thành phẩm lên đúng xe/container để giao cho khách hàng.

### 8.2. Chuẩn bị màn hình

1. Mở trang **Scan Ship KTP** và đăng nhập.
2. Nhập **Số container**.
3. Nhập **Biển số xe**.
4. Bấm chuột vào ô quét mã. Ô quét được sử dụng ngay sau khi đã nhập đủ container và biển số; không cần chờ cân mới được focus vào ô.
5. Kiểm tra cân đã tự kết nối, số cân lớn hơn `0` và ổn định trước khi gửi mã quét.
6. Màn hình chỉ hiển thị đồng hồ ngày giờ hiện tại; người vận hành không chọn ngày xuất.

Sau carton đầu tiên thành công, số container và biển số bị khóa để tránh đổi xe giữa phiên. Muốn đổi xe/container phải mở **Phiên mới**.

### 8.3. Các bước quét

1. Đặt carton lên cân và chờ số cân ổn định.
2. Quét `CARTONBAR` hoặc LPN.
3. Hệ thống kiểm tra trạng thái carton, Shipping Plan và trọng lượng.
4. Nếu carton khác RY với danh sách phiên, đọc cảnh báo; nếu đúng carton cần xuất thì hệ thống vẫn tiếp tục.
5. Nếu đạt điều kiện, hệ thống xuất ngay một carton; không cần bấm Upload.
6. Xác nhận thông báo thành công; kiểm tra số thùng, số đôi, `rgw` và `ogw` trên danh sách.
7. Chọn carton để xem chi tiết. Khu vực chi tiết chỉ hiển thị một dòng đúng size/XH của carton đó.
8. Nhấc carton khỏi cân, chờ cân sẵn sàng và thực hiện carton tiếp theo.

### 8.4. Điều kiện Scan Ship

#### A. Carton và trạng thái

- Tìm thấy carton trong `YWCP`.
- Carton có RY.
- Bắt buộc `SB=1`.
- `EXEDATE` phải rỗng.

#### B. Container, biển số và Shipping Plan

- Container và biển số không được để trống.
- Nếu carton chưa có dòng trong `ShippingPlanCarton`, hệ thống dùng container và biển số người vận hành đã nhập.
- Nếu carton có Shipping Plan, cả container và biển số phải trùng với kế hoạch; so sánh không phân biệt chữ hoa/thường.
- Sai một trong hai thông tin thì không được Scan Ship.

#### C. Cân nặng

- Cân phải kết nối và trọng lượng xuất `ogw > 0`.
- Ưu tiên dùng trọng lượng Scan In `rgw` làm tham chiếu, sai số cho phép `±0,05 kg`.
- Nếu không có `rgw > 0`, dùng trọng lượng chuẩn `sgw`, sai số cho phép `±0,5 kg`.
- Nếu cả `rgw` và `sgw` không hợp lệ, carton bị từ chối và không thể xác nhận bỏ qua.

### 8.5. Xử lý lệch cân Scan Ship

1. Khi trọng lượng vượt giới hạn, kiểm tra lại đúng carton, cân và hàng bên trong.
2. Chọn **Hủy** để dừng; carton không được xuất.
3. Nếu được phép xuất, chọn **Xác nhận**.
4. Hệ thống cập nhật xuất hàng và ghi log xác nhận lệch cân vào cùng bảng với Scan In: `LYS_WEB.dbo.KTP_Scan_logs`.
5. Log gồm carton/RY, người xác nhận, trọng lượng tham chiếu, trọng lượng xuất, sai số, nguồn tham chiếu `RGW/SGW`, container và biển số.

### 8.6. Database được cập nhật khi Scan Ship thành công

| Bảng/field | Giá trị cập nhật |
|---|---|
| `YWCP.SB` | `3` |
| `YWCP.EXEDATE` | Thời gian xuất hiện tại |
| `YWCP.CON_NO` | Số container đã nhập |
| `YWCP.BIEN_SO` | Biển số xe đã nhập |
| `YWCP.ogw` | Trọng lượng thực tế khi xuất |
| `LYS_WEB.dbo.KTP_Scan_logs` | Chỉ thêm log khi có xác nhận lệch cân |

### 8.7. Kết thúc/Refresh phiên Scan Ship

1. Đối chiếu tổng số thùng và số đôi với hàng thực tế trên xe/container.
2. Kiểm tra lại container và biển số của phiên.
3. Bấm **Phiên mới** sau khi hoàn tất xe/container hiện tại.
4. Hệ thống xóa danh sách phiên trên màn hình đồng thời xóa container và biển số để nhập chuyến mới.
5. Thao tác này không hoàn tác các carton đã Scan Ship trong ERP.

## 9. Xử lý lỗi và bất thường

| Hiện tượng | Xử lý của người vận hành |
|---|---|
| Phiên đăng nhập hết hạn/`401 Unauthorized` | Dừng quét, đăng nhập lại tại đúng trang Scan In/Out/Ship, sau đó kiểm tra trạng thái carton cuối trước khi quét lại. |
| Không có quyền/`403 Forbidden` | Dùng đúng tài khoản được cấp quyền hoặc liên hệ quản lý/IT; không mượn tài khoản. |
| Không tìm thấy carton | Lau tem, quét lại đúng mã; nếu vẫn lỗi, nhờ ERP kiểm tra `CARTONBAR/LPNBarCode`. |
| Carton khác RY | Đây là cảnh báo, không phải lỗi khóa. Kiểm tra đúng thùng cần thao tác rồi tiếp tục. |
| Carton đã có trong phiên | Không quét lại; kiểm tra dòng đã thành công trong danh sách. |
| Sai `SB` | Không sửa tay. Xác định carton đang ở bước nào theo bảng trạng thái và thực hiện đúng bước trước đó. |
| Cân mất kết nối hoặc không có số cân | Dừng Scan In/Ship, kiểm tra cáp/nguồn/cổng cân và kết nối lại. Không nhập trọng lượng thủ công. |
| Lệch cân | Cân lại; chỉ xác nhận bỏ qua sau khi kiểm tra thực tế và được phép. |
| Sai container/biển số so với Shipping Plan | Không xuất; kiểm tra lại kế hoạch hoặc nhờ bộ phận phụ trách điều chỉnh dữ liệu. |
| API/database lỗi hoặc timeout | Không mặc định là chưa ghi. Giữ nguyên carton, kiểm tra trạng thái ERP hoặc nhờ IT đối chiếu trước khi quét lại. |
| Không ghi được log lệch cân | Xem thao tác là thất bại; không tự ý di chuyển carton cho đến khi kiểm tra lại trạng thái ERP và log. |

## 10. Kiểm soát cuối ca

1. Đối chiếu số thùng và số đôi của từng phiên với hàng thực tế.
2. Kiểm tra các carton báo lỗi hoặc bị hủy xác nhận lệch cân.
3. Đảm bảo không còn carton đã di chuyển vật lý nhưng màn hình chưa báo thành công.
4. Với Scan Ship, đối chiếu `CON_NO`, `BIEN_SO`, số thùng và số đôi theo từng xe/container.
5. Báo quản lý/IT danh sách carton cần đối chiếu, kèm mã carton, RY, chức năng, thời gian và nội dung lỗi.

## 11. Dữ liệu truy vết chính

- `LYS_ERP.dbo.YWCP`: trạng thái và thông tin Scan In/Out/Ship của từng carton.
- `LYS_WEB.dbo.KTP_Scan_logs`: log xác nhận lệch cân của Scan In và Scan Ship.
- Danh sách phiên trên trình duyệt: hỗ trợ vận hành tại chỗ, không thay thế dữ liệu ERP và không phải chức năng hoàn tác.

## 12. Checklist nhanh cho người vận hành

### Scan In

- [ ] Đúng chế độ Carton hoặc Carton and Inbox.
- [ ] Đúng Production Line, Storage và Location.
- [ ] Cân kết nối, số cân ổn định.
- [ ] Carton đang ở `SB` rỗng, `2` hoặc `4`.
- [ ] Màn hình báo thành công trước khi chuyển thùng.

### Scan Out

- [ ] Đã chọn đúng Tái chế hoặc Kiểm hàng.
- [ ] Carton đang ở `SB=1`.
- [ ] Màn hình báo thành công và bộ đếm tăng đúng.

### Scan Ship

- [ ] Đúng container và biển số.
- [ ] Cân kết nối, số cân ổn định.
- [ ] Carton đang ở `SB=1` và chưa có `EXEDATE`.
- [ ] Đúng Shipping Plan nếu carton có kế hoạch.
- [ ] Màn hình báo thành công trước khi đưa thùng lên xe.
