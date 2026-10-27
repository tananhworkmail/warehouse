# Tổng hợp điều kiện Scan In – Kho thành phẩm (KTP)

Tài liệu này mô tả **đúng logic web hiện đang thực thi**, gồm bước kiểm tra trước khi scan, điều kiện cho phép nhập kho và dữ liệu được ghi vào ERP. Đây không phải danh sách đề xuất hoặc bản mô tả nguyên trạng của ERP cũ.

## 1. Quyền truy cập

Toàn bộ API Scan In đều yêu cầu:

- JWT đăng nhập hợp lệ.
- Có bản ghi quyền trong `Blimits`

## 2. Thông tin phiên trước khi scan

Web hiện dùng:

| Thông tin | Quy tắc hiện tại |
|---|---|
| Production line | Bắt buộc chọn |
| Storage (`KCBH`) | Bắt buộc; mặc định `A2`, được phép sửa |
| Location (`KVBH`) | Không bắt buộc; được phép để trống |

Khi submit, backend kiểm tra thêm:

- Warehouse `A` phải tồn tại trong `KCCKCP.CKBH`.
- Production line phải tồn tại trong `BDepartment` và thỏa `GSBH = 'VDH'`, `UserID = '340205'`, `ProYN = 1`, `YN = '1'`, `GXLB = 'A'`.

Phiên làm việc được lưu theo userid trong `localStorage`, nên reload trang vẫn giữ Production line, Storage, Location, chế độ scan và danh sách carton. Chỉ nút **Kết thúc phiên** mới xóa phiên; thao tác này reset Production line và Location, còn Storage giữ giá trị hiện tại/mặc định `A2`.

## 3. Nhận diện carton

Web nhận cả `CARTONBAR` và `LPNBarCode`:

1. Tìm chính xác `YWCP.CARTONBAR` theo mã vừa quét.
2. Nếu không thấy, tìm `YWCP.LPNBarCode`.
3. Nếu chuỗi LPN dài hơn 15 ký tự, chỉ lấy 15 ký tự cuối để tìm.
4. Không tìm thấy carton thì chặn Scan In.

Thông tin carton được ghép từ `YWCP`, `YWDD`, `DDZL`, `YWBZPOS`, `XXZL`, `BDepartment`, `LBZLS`, `Bgszl` và `YWXXNew`.

## 4. Điều kiện carton và trạng thái `SB`

Carton chỉ hợp lệ khi:

- Có đủ `DDBH` và đơn hàng `OrderNo/YSBH`.
- Chưa tồn tại trong danh sách carton của phiên hiện tại.
- Cùng đơn hàng với các carton đã có trong phiên.
- `SB` thuộc một trong các trạng thái được phép.

| Giá trị `SB` | Ý nghĩa xử lý trên web |
|---|---|
| `NULL` hoặc rỗng | Thành phẩm mới, cho phép nhập KTP lần đầu |
| `2` | Luồng tái nhập/repacking, cho phép Scan In |
| `4` | Luồng tái nhập sau kiểm tra chất lượng, cho phép Scan In |
| `1` | Đang ở trong kho thành phẩm, chặn |
| `3` | Đã xuất khỏi kho thành phẩm, chặn |
| Giá trị khác | Chặn |

## 5. Kiểm tra cân nặng

### 5.1 Dữ liệu chuẩn

- Trọng lượng chuẩn carton: `YWCP.sgw`; phải lớn hơn `0`.
- Trọng lượng thực tế: giá trị `G.W.` đọc từ cân; phải lớn hơn `0`.
- Giá trị `Gram`:
  - Nhà máy `VB1` hoặc `VB7`: `ROUND(sgw * 0.03, 3)`.
  - Các nhà máy khác: `YWXXNew.NWeigh * 1.04`, tra theo model, màu và size.
- `Gram` phải lớn hơn `0`.

### 5.2 Công thức web đang bắt

```text
tolerance = Gram / 3
difference = sgw - actualWeight
đạt khi ABS(difference) <= tolerance
```

Web đang dùng so sánh trị tuyệt đối cho **tất cả nhà máy**. Thiếu `sgw`, thiếu `Gram` hoặc không đọc được `G.W.` đều không đạt.

### 5.3 Mở khóa sai cân

Nút mở khóa chỉ hiện khi:

- Riêng điều kiện cân nặng đang không đạt.
- Điều kiện carton, sản lượng và vật tư đều đạt.

Khi web nhận thấy chỉ riêng cân nặng không đạt, hộp thoại sau tự động xuất hiện:

> Trọng lượng thùng hàng không đạt chuẩn. Có muốn xác nhận tiếp tục nhập hàng?

- Chọn **Không**: hủy xác nhận, không ghi carton và reset số cân.
- Chọn **Có**: mở ô `Nhập xác nhận cho phép lệch cân`.

Ô xác nhận được hiển thị dưới dạng `password` để che ký tự và giao diện không gợi ý giá trị phải nhập. Giá trị hợp lệ theo nghiệp vụ hiện tại vẫn là userid đang đăng nhập. Backend không tin dữ liệu phía trình duyệt mà kiểm tra `weightOverrideUserId` phải trùng với userid lấy từ JWT.

Mở khóa chỉ chuyển riêng kiểm tra cân sang đạt; các điều kiện carton, sản lượng, vật tư và hộp con vẫn phải đạt. Backend chỉ chấp nhận mở khóa khi kết quả cân ban đầu thực sự là `failed`.

Web không ghi sự kiện mở khóa cân vào file `system.log`. Khi carton lệch cân đã UPDATE và được kiểm tra lại thành công, service ghi audit vào `LYS_WEB.dbo.KTP_Scan_logs` trên máy `.21`, gồm userid người cho phép, userid thực hiện scan, carton, đơn hàng, snapshot cân và thông tin vị trí nhập kho. Nếu INSERT log thất bại, transaction cập nhật `YWCP` trả lỗi và rollback. `YWCP.INUser` vẫn lưu userid thực hiện nhập kho như luồng ERP.

## 6. Kiểm tra sản lượng Gò theo size

Điều kiện này chỉ khóa với carton thuộc nhà máy có `Bgszl.SFL = 'PD'`. Nhà máy không thuộc nhóm `PD` được cho qua điều kiện sản lượng theo đúng nhánh web hiện tại.

Web tính:

```text
completedQty = SUM(SMDDSS.Qty * SMDDSS.okCTS)
```

với:

- `SMDDSS.DDBH LIKE OrderNo + '%'`
- `SMDDSS.GXLB = 'A'`
- `SMDDSS.XXCC = size của carton`

Số lượng sau scan:

```text
afterScanQty = alreadyInQty + sessionQty + currentCartonQty
```

Trong đó:

- `alreadyInQty`: tổng `YWBZPOS.Qty` (fallback `YWCP.Qty`) của cùng đơn hàng và size, có `LastInDate`, `SB = '1'`, không tính carton hiện tại và các carton trong phiên.
- `sessionQty`: tổng số lượng các carton cùng đơn hàng và size đã có trong phiên.
- `currentCartonQty`: `YWBZPOS.Qty` (fallback `YWCP.Qty`) của carton đang scan.

Điều kiện đạt:

- `completedQty > 0`.
- `afterScanQty <= completedQty`.

Ngoại lệ: nếu `YWDD.Flag = '1'`, web cho qua kiểm tra sản lượng Gò.

## 7. Kiểm tra cấp phát vật tư

### 7.1 Xác định đơn vật tư

Web ưu tiên lấy đơn gốc từ `DDZLTR.DDBH` với `DDZLTR.DDBH1 = OrderNo`; không có thì dùng chính `OrderNo`.

### 7.2 Nguồn định mức

Web tổng hợp các dòng vật tư cần dùng từ:

- Vật tư chính: `DDZL_ZLZLS2`.
- Gia công (`SFL = 'GC'`): `ZLZLS2_YM`.
- Cao su (`SFL = 'RB'`): `LIY_DD.dbo.ZLZLS3`.

Phạm vi nhà máy được lấy theo cùng `DFL` trong `Bgszl` và bỏ các dòng `tybh = 'Y'`. Nếu đơn hàng cần nhánh cao su nhưng không đọc được `LIY_DD.dbo.ZLZLS3`, điều kiện vật tư trả `NO_DATA` và chặn Scan In.

### 7.3 Số lượng đã cấp và mức bao phủ

Vật tư đã cấp lấy từ `KCLLS`, nối `KCLL`, chỉ tính phiếu có `KCLL.CFMID <> 'NO'`.

Cho từng mã vật tư:

```text
coveredPairs = issuedQty * orderPairs / requiredQty
```

Mức bao phủ của đơn là `coveredPairs` nhỏ nhất trong tất cả vật tư bắt buộc.

Điều kiện dữ liệu tối thiểu:

- Có đơn vật tư.
- `orderPairs > 0`.
- Có ít nhất một dòng định mức.
- Không có dòng vật tư bắt buộc nào có số lượng cấp `<= 0`.

### 7.4 Quy tắc đạt

Với carton nhập lần đầu (`SB` rỗng), đạt khi một trong các nhánh sau đúng:

- `coveredPairs >= orderPairs`; hoặc
- `alreadyInQty + sessionQty + carton.Qty < coveredPairs`; hoặc
- `YWDD.Flag = '1'`.

Nếu không đúng thì báo vật tư cấp phát chưa đủ cho carton hiện tại.

Với carton tái nhập (`SB = '2'` hoặc `SB = '4'`), web vẫn yêu cầu đọc được định mức và không có dòng vật tư hoàn toàn chưa mở phiếu cấp phát; sau đó không khóa lại theo mức số lượng bao phủ.

Lưu ý: `YWDD.Flag = '1'` không bỏ qua lỗi thiếu dữ liệu định mức hoặc có dòng vật tư chưa mở phiếu cấp phát, vì các lỗi này được kiểm tra trước nhánh Flag.

## 8. Chế độ Scan Carton and Inbox

Chế độ `carton` không bắt điều kiện hộp con. Chế độ `carton_inbox` bắt buộc kiểm tra và quét đủ hộp.

Mã hộp được lấy theo thứ tự:

1. `YWBZPOS.InBoxCode` của đúng `DDBH/XH`.
2. Nếu trống, lấy `YWCPLH.InBoxBar` mới nhất theo model, màu, size và quốc gia của đơn; cho phép mapping quốc gia cụ thể hoặc `NORMAL`, ưu tiên quốc gia cụ thể.

Số hộp phải quét lấy từ `YWBZPOS.Qty`. Các dòng có cùng mã hộp được cộng dồn.

Chế độ này chỉ đạt khi:

- Có mã hộp cho tất cả dòng size có số lượng lớn hơn `0`.
- Tổng số hộp dự kiến lớn hơn `0`.
- Mỗi mã quét thuộc carton.
- Không mã nào được quét vượt số lượng dự kiến.
- Tổng số lần quét của từng mã đúng bằng số lượng dự kiến.

Quét thiếu, quét dư hoặc quét mã hộp không thuộc carton đều chặn ghi carton.

## 9. Điều kiện quyết định cuối cùng

Backend chỉ cho ghi khi:

```text
canScanIn = eligibility.passed
         AND weightRule.passed
         AND productionCheck.passed
         AND materialCheck.passed
         AND inboxCheck.passed
```

Ở chế độ `carton`, `inboxCheck` tự động đạt. Ở chế độ `carton_inbox`, `inboxCheck` chỉ đạt khi đã quét đủ và đúng hộp.

Web tự động xử lý sau khi máy quét gửi mã:

- Carton thường: đủ điều kiện thì ghi ngay, không cần nhấn Enter lần nữa.
- Carton and Inbox: carton đạt các điều kiện chính thì chuyển sang thu mã hộp; đủ hộp thì ghi ngay.
- Không đạt: phát âm báo lỗi, xóa ô mã đang nhập và reset số cân để người dùng scan lại.
- Thành công: phát âm báo thành công, chạy hiệu ứng carton xuống bảng, giữ thông tin preview, thêm dòng vào bảng phiên và reset số cân cho carton tiếp theo.

## 10. Kiểm tra lại trong transaction trước khi ghi

Kết quả kiểm tra lúc xem preview không được dùng để ghi thẳng. Khi submit, backend:

1. Mở transaction.
2. Đọc lại carton.
3. Khóa nghiệp vụ bằng `sp_getapplock`, ưu tiên khóa theo OrderNo, sau đó DDBH, rồi CARTONBAR.
4. Đọc lại carton sau khi có khóa.
5. Kiểm tra lại Warehouse và Production line.
6. Chạy lại tuần tự toàn bộ điều kiện carton, cân, sản lượng, vật tư và hộp con.
7. Áp dụng mở khóa cân nếu yêu cầu hợp lệ.
8. Chỉ UPDATE khi `SB IS NULL OR SB IN ('', '2', '4')`.
9. Đọc lại dòng vừa ghi và đối chiếu các field quan trọng.

Bất kỳ bước nào lỗi thì transaction rollback; không để lại bản ghi nhập kho dở dang.

## 11. Dữ liệu được ghi vào `YWCP`

Khi thành công, backend cập nhật:

| Field | Giá trị/quy tắc |
|---|---|
| `SB` | `'1'` |
| `CKBH` | `'A'` |
| `DepNo` | Production line nếu là lần nhập đầu (`INDATE IS NULL`); tái nhập giữ `DepNo` cũ |
| `KCBH` | Storage người dùng chọn |
| `KVBH` | Location; có thể rỗng |
| `INCS` | `1` nếu chưa từng nhập; nếu đã có `INDATE` thì tăng thêm `1` |
| `INDATE` | `GETDATE()` ở lần nhập đầu; tái nhập giữ ngày cũ |
| `LastInDate` | `GETDATE()` |
| `rgw` | Trọng lượng thực tế đã dùng để kiểm tra |
| `INUser` | Userid lấy từ JWT |
| `REDATE` | `GETDATE()` nếu trạng thái trước đó là `2` hoặc `4` |
| `RECS` | Tăng `1` nếu trạng thái trước đó là `2` hoặc `4` |

Sau UPDATE, backend xác nhận lại `CARTONBAR`, `SB`, `CKBH`, `DepNo`, `KCBH`, `KVBH`, `INCS`, `INDATE`, `LastInDate`, `rgw`, `INUser`; riêng tái nhập còn xác nhận `REDATE` và `RECS`.

## 12. Tóm tắt các trường hợp bị chặn

| Nhóm | Trường hợp bị chặn |
|---|---|
| Quyền | Chưa đăng nhập hoặc không có quyền `VDH/21/SN7` |
| Thiết lập | Thiếu Production line/Storage; Warehouse A hoặc Production line không hợp lệ |
| Carton | Không tìm thấy barcode/LPN; thiếu đơn; trùng phiên; khác đơn trong phiên |
| Trạng thái | `SB = 1`, `SB = 3` hoặc trạng thái ngoài rỗng/2/4 |
| Cân | Thiếu `sgw`, thiếu `Gram`, chưa có G.W. hoặc lệch vượt `Gram/3` mà chưa mở khóa |
| Sản lượng | Nhóm PD chưa có sản lượng Gò hoặc số lượng sau scan vượt sản lượng hoàn thành |
| Vật tư | Thiếu định mức, thiếu dữ liệu cao su, có vật tư chưa mở cấp phát hoặc mức cấp không đủ |
| Hộp con | Thiếu mapping, quét sai mã, quét dư hoặc chưa đủ số hộp trong chế độ Carton and Inbox |
| Ghi dữ liệu | Carton đổi trạng thái trong lúc xử lý, UPDATE không đúng một dòng hoặc kiểm tra sau ghi không khớp |
