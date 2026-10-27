unit ScanQRDelivery1;

interface

uses
  Windows, Messages, SysUtils, Variants, Classes, Graphics, Controls, Forms,
  Dialogs, DBTables, DB, ScktComp, Menus, GridsEh, DBGridEh, ExtCtrls,pngimage,
  StdCtrls, Mask, DBCtrls, Buttons, ComCtrls, comobj, jpeg, Spin, DelphiZXIngQRCode,
  QRCtrls, QuickRpt,ShellApi, Printers, RpCon, RpConDS, RpConBDE, RpDefine, RpRave,
  GIFImage;

type
  TScanQRDelivery = class(TForm)
    Panel3: TPanel;
    Image2: TImage;
    Label2: TLabel;
    imgClose00: TImage;
    LLBarCodeEdit: TEdit;
    Edit1: TEdit;
    Panel1: TPanel;
    Panel2: TPanel;
    Panel4: TPanel;
    Panel5: TPanel;
    Image1: TImage;
    Image3: TImage;
    Panel6: TPanel;
    QRImage1: TImage;
    KCLLSQry: TQuery;
    KCLLSQryCLBH: TStringField;
    KCLLSQryywpm: TStringField;
    KCLLSQryDWBH: TStringField;
    KCLLSQryTempQty: TCurrencyField;
    KCLLSQryQty: TCurrencyField;
    KCLLSQryYWSM: TStringField;
    KCLLSQryBLSB: TStringField;
    KCLLSQryLLNO: TStringField;
    KCLLSQrySCBH: TStringField;
    KCLLSQryDFL: TStringField;
    KCLLSQrySCANUSERID: TStringField;
    KCLLSQrySCANUSERDATE: TDateTimeField;
    DS2: TDataSource;
    DBGridEh2: TDBGridEh;
    Report: TImage;
    Panel7: TPanel;
    QRImage2: TImage;
    Report1: TImage;
    imgSave00: TImage;
    LLBarCodeEdit1: TEdit;
    GroupBox2: TGroupBox;
    Label1: TLabel;
    Label7: TLabel;
    Label8: TLabel;
    Label9: TLabel;
    Label12: TLabel;
    Label14: TLabel;
    Label15: TLabel;
    Label16: TLabel;
    Label17: TLabel;
    Label18: TLabel;
    DBMemo1: TDBMemo;
    DBMemo2: TDBMemo;
    DBMemo3: TDBMemo;
    DBMemo4: TDBMemo;
    DBEdit1: TDBEdit;
    DBEdit3: TDBEdit;
    DBEdit4: TDBEdit;
    DBEdit5: TDBEdit;
    DBEdit7: TDBEdit;
    Image4: TImage;
    DataSource1: TDataSource;
    Query1: TQuery;
    Query1SCNO: TStringField;
    Query1CLBH: TStringField;
    Query1YWPM: TStringField;
    Query1ZSYWJC: TStringField;
    Query1Memo_RY: TStringField;
    Query1Memo_Article: TStringField;
    Query1Qty: TFloatField;
    Query1PACK1: TStringField;
    Query1Pack: TIntegerField;
    Query1Memo: TStringField;
    Query1DWBH: TStringField;
    Query1CGNO: TStringField;
    Query1QRCode: TStringField;
    Query1FIFO: TStringField;
    Query1LotNO: TStringField;
    tmpQry: TQuery;
    CLBH: TDBText;
    YWPM: TDBText;
    LLNO: TDBText;
    Track: TImage;
    Save: TImage;
    KCLLSQryAllqty: TCurrencyField;
    KCLLSQryLLQty: TCurrencyField;
    KCLLSQrySKU: TStringField;
    Query1IDCheck: TStringField;
    UpdateSQL1: TUpdateSQL;
    ImageOUT: TImage;
    lbOUT: TLabel;
    KCLLSQrycheckry: TFloatField;
    procedure FormDestroy(Sender: TObject);
    procedure FormClose(Sender: TObject; var Action: TCloseAction);
    procedure imgClose00Click(Sender: TObject);
    procedure FormCreate(Sender: TObject);
    procedure ReportClick(Sender: TObject);
    procedure Report1Click(Sender: TObject);
    procedure LLBarCodeEdit1KeyPress(Sender: TObject; var Key: Char);
    procedure LLBarCodeEditKeyPress(Sender: TObject; var Key: Char);
    procedure Edit1KeyPress(Sender: TObject; var Key: Char);
    procedure imgSave00Click(Sender: TObject);
    procedure KCLLSQryAfterOpen(DataSet: TDataSet);
    procedure TrackClick(Sender: TObject);
    procedure SaveClick(Sender: TObject);
    procedure DBGridEh2GetCellParams(Sender: TObject; Column: TColumnEh;
      AFont: TFont; var Background: TColor; State: TGridDrawState);
    procedure DBGridEh2CellClick(Column: TColumnEh);
    procedure ReportDblClick(Sender: TObject);
    procedure DBGridEh2DblClick(Sender: TObject);
    procedure ImageOUTClick(Sender: TObject);
  private
    QRCodeBitmap:TBitmap;
    LLNOlist:string;
    { Private declarations }
  public
    QRCodeBitmap1:TBitmap;
    { Public declarations }
    procedure UpLoadQRImg(QRdata:string);
    procedure UpLoadQRImg1(QRdata:string);
  end;

var
  ScanQRDelivery: TScanQRDelivery;

implementation

uses main1, ScanMaterialQty1, ScanQRDeliveryOWE1, ScanQR1, ScanQRDeliveryPAY1, FunUnit;

{$R *.dfm}

procedure TScanQRDelivery.FormDestroy(Sender: TObject);
begin
  ScanQRDelivery:=nil;
end;

procedure TScanQRDelivery.FormClose(Sender: TObject;
  var Action: TCloseAction);
begin
  Action:=Cafree;
end;

procedure TScanQRDelivery.imgClose00Click(Sender: TObject);
begin
  if MessageDlg('Ban co muon dong khong?',mtConfirmation, [mbYes, mbNo], 0) = mrYes then
  begin
    Close;
  end;
end;

procedure TScanQRDelivery.FormCreate(Sender: TObject);
var imout, imscan, imscan1, imprint, imtrack, imsave: TPNGObject;
    GIF1, GIF2: TGIFimage;
begin
  imscan := TPNGObject.Create;
  imscan.LoadFromFile('btnscan.png');
  Report.Picture.Assign(imscan);

  imscan1 := TPNGObject.Create;
  imscan1.LoadFromFile('btnscan.png');
  Report1.Picture.Assign(imscan1);

  imprint := TPNGObject.Create;
  imprint.LoadFromFile('btnprint.png');
  imgSave00.Picture.Assign(imprint);

  imtrack := TPNGObject.Create;
  imtrack.LoadFromFile('btntrack.png');
  Track.Picture.Assign(imtrack);

  imout := TPNGObject.Create;
  imout.LoadFromFile('btnlogout.png');
  imgClose00.Picture.Assign(imout);

  imsave := TPNGObject.Create;
  imsave.LoadFromFile('btnsave.png');
  Save.Picture.Assign(imsave);

  GIF1 := TGIFImage.Create;
  GIF1.LoadFromFile('scan.gif');
  QRImage1.Picture.Assign(GIF1);

  GIF2 := TGIFImage.Create;
  GIF2.LoadFromFile('scan.gif');
  QRImage2.Picture.Assign(GIF2);

  imgClose00.Width:=Round(Screen.Width/28);
  imgClose00.Height:=Round(Screen.Width/28);
  imgClose00.Left:=Round(Screen.Width-imgClose00.Width);
  imgClose00.Top:=Round(Screen.DesktopTop+5);

  Track.Width:=imgClose00.Width;
  Track.Height:=imgClose00.Height;
  Track.Left:=imgClose00.Left-Track.Width-imgClose00.Width;
  Track.Top:=imgClose00.Top;

  Panel1.Width:=Round(Screen.Width-GroupBox2.Width-10);

  YWPM.Width:=Round(Panel1.Width/2);

  imscan.free;
  imscan1.free;
  imprint.free;
  imout.free;
  imtrack.free;
  imsave.free;
  GIF1.free;
  GIF2.free;
end;

procedure TScanQRDelivery.ReportClick(Sender: TObject);
  var list: string;
begin
  LLBarCodeEdit1.Visible:=false;
  LLBarCodeEdit.Visible:=true;
  LLBarCodeEdit.SetFocus;
  edit1.ReadOnly:=false;
  KCLLSQry.active:=false;
  KCLLSQry.active:=true;
end;

procedure TScanQRDelivery.Report1Click(Sender: TObject);
begin
  LLBarCodeEdit.Visible:=false;
  LLBarCodeEdit1.Visible:=true;
  LLBarCodeEdit1.SetFocus;
  imgSave00.Visible:=false;
  Save.Visible:=false;
  Query1.active:=true;
  imgSave00.Visible:=true;
end;

procedure TScanQRDelivery.Edit1KeyPress(Sender: TObject; var Key: Char);
  var s:string;
begin
  s := ('1234567890.'#8); //Add chars you want to allow
  if pos(key,s) =0 then Key:=#0;
end;

procedure TScanQRDelivery.UpLoadQRImg(QRdata:string);
var
  QRCode: TDelphiZXingQRCode;
  Row,Column: Integer;
begin
  QRCode := TDelphiZXingQRCode.Create;
  try
    QRCode.Data:=QRdata;
    QRCode.Encoding := TQRCodeEncoding(3);
    QRCode.QuietZone := StrToIntDef('1',4);
    QRCodeBitmap.Width:=QRCode.Rows;
    QRCodeBitmap.Height:=QRCode.Columns;
    for Row := 0 to QRCode.Rows - 1 do
    begin
      for Column := 0 to QRCode.Columns - 1 do
      begin
        if (QRCode.IsBlack[Row, Column]) then
        begin
          QRCodeBitmap.Canvas.Pixels[Column, Row] := clBlack;
        end else
        begin
          QRCodeBitmap.Canvas.Pixels[Column, Row] := clWhite;
        end;
      end;
    end;
  finally
    QRCode.Free;
  end;
end;

procedure TScanQRDelivery.UpLoadQRImg1(QRdata:string);
var
  QRCode: TDelphiZXingQRCode;
  Row,Column: Integer;
begin
  QRCode := TDelphiZXingQRCode.Create;
  try
    QRCode.Data:=QRdata;
    QRCode.Encoding := TQRCodeEncoding(3);
    QRCode.QuietZone := StrToIntDef('1',4);
    QRCodeBitmap1.Width:=QRCode.Rows;
    QRCodeBitmap1.Height:=QRCode.Columns;
    for Row := 0 to QRCode.Rows - 1 do
    begin
      for Column := 0 to QRCode.Columns - 1 do
      begin
        if (QRCode.IsBlack[Row, Column]) then
        begin
          QRCodeBitmap1.Canvas.Pixels[Column, Row] := clBlack;
        end else
        begin
          QRCodeBitmap1.Canvas.Pixels[Column, Row] := clWhite;
        end;
      end;
    end;
  finally
    QRCode.Free;
  end;
end;

procedure TScanQRDelivery.LLBarCodeEdit1KeyPress(Sender: TObject;
  var Key: Char);
var barcode, qty, remqty, SCANQty, CLBH1, DFL, SCBH, CFMDATE1, y,m, a: string;
    i:integer;
begin
  if (key='.') or (key=#13) then
  begin
    if edit1.Text='' then edit1.Text:='0';
    QRCodeBitmap1.Free;
    barcode:=LLBarCodeEdit1.Text;
    with tmpQry do
    begin
      Active:=false;
      SQL.Clear;
      SQL.Add('select isnull(KCRKScan_RFSS.RemQty,KCRKScan_RFSS.Qty) RemQty, KCRKScan_RFSS.CLBH, CFMDATE from KCRKScan_RFSS ');
      SQL.Add('Where KCRKScan_RFSS.barcode='''+barcode+''' ');
      active:=true;
      remqty:=tmpQry.fieldbyname('RemQty').AsString;
      CLBH1:=tmpQry.fieldbyname('CLBH').AsString;
      CFMDATE1:=tmpQry.fieldbyname('CFMDATE').AsString;
    end;
    if CFMDATE1='' then
    begin
      Showmessage('Vui long scan nhap kho !!!');
      edit1.Text:='0';
      KCLLSQry.Active:=false;
      KCLLSQry.cachedupdates:=false;
      KCLLSQry.requestlive:=false;
      KCLLSQry.Active:=true;
      LLBarCodeEdit1.Text:='';
      key:=#13;
      Exit;
    end;
    if (strtofloat(edit1.Text)>=0) and (edit1.ReadOnly=false) then
    begin
      qty:=FloatToStr(StrToFloat(remqty)-StrToFloat(edit1.Text));
      with tmpQry do          //取服務器的年月值
      begin
        active:=false;
        sql.Clear;
        sql.Add('select year(getdate()) as FY,month(getdate()) as FM ');
        active:=true;
        y:=fieldbyname('FY').asstring;
        m:=fieldbyname('FM').asstring;
        active:=false;
      end;
      if length(m)<2 then
      m:='0'+m;
      with tmpQry do    //計算領料單流水號
      begin
        active:=false;
        sql.Clear;
        sql.Add('select LLNO_OWE from KCLLScan_RFSSS where LLNO_OWE like '+''''+y+m+'%'+'''');
        sql.add('order by LLNO_OWE');
        active:=true;
        if recordcount >0 then
          begin
            last;
            a:=fieldbyname('LLNO_OWE').AsString;
            a:=copy(a,7,5);
            a:=inttostr(strtoint(a)+1);
            while length(a)<5 do
              begin
                a:='0'+a;
              end;
          end
            else
              begin
                a:='00001';
              end;
        a :=y+m+a;
        active:=false;
      end;
      if StrToFloat(remqty)>0 then
      begin
        if StrToFloat(edit1.Text)=0 then
        begin
          with tmpQry do
          begin
            Active:=false;
            SQL.Clear;
            sql.Add('Insert into KCLLScan_RFSSS (LLNO_OWE, LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE) ');
            SQL.Add('Values( '''+a+''',''NO'','''+CLBH1+''',convert(varchar(15), Getdate(), 114),convert(varchar(15), Getdate(), 111), ');
            SQL.Add('        '''+barcode+''','''+remqty+''','''+main.Edit1.Text+''',Getdate()) ');
            ExecSQL;
            Active:=false;
            SQL.Clear;
            sql.Add('Update KCRKScan_RFSS set RemQty=0, CFMID1='''+main.Edit1.Text+''', CFMDATE1=getdate() ');
            SQL.Add('Where KCRKScan_RFSS.barcode='''+barcode+''' and CFMDATE is not null ');
            ExecSQL;
          end;
        end else if (StrToFloat(edit1.Text)>0) and (StrToFloat(qty)>=0) then begin
          with tmpQry do
          begin
            Active:=false;
            SQL.Clear;
            sql.Add('Insert into KCLLScan_RFSSS (LLNO_OWE, LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE) ');
            SQL.Add('Values( '''+a+''',''NO'','''+CLBH1+''',convert(varchar(15), Getdate(), 114),convert(varchar(15), Getdate(), 111), ');
            SQL.Add('        '''+barcode+''','''+edit1.Text+''','''+main.Edit1.Text+''',Getdate()) ');
            ExecSQL;
            Active:=false;
            SQL.Clear;
            sql.Add('Update KCRKScan_RFSS set RemQty='''+qty+''', CFMID1='''+main.Edit1.Text+''', CFMDATE1=getdate() ');
            SQL.Add('Where KCRKScan_RFSS.barcode='''+barcode+''' and CFMDATE is not null ');
            ExecSQL;
          end;
        end;
      end;
    end else if (edit1.ReadOnly=true) then
    begin
      try
        KCLLSQry.first;
        for i:=1 to KCLLSQry.RecordCount do
        begin
          if (CLBH1=KCLLSQry.fieldbyname('CLBH').AsString) and (KCLLSQry.fieldbyname('checkry').Value=1) then
          begin
            DFL:=KCLLSQry.fieldbyname('DFL').AsString;
            SCBH:=KCLLSQry.fieldbyname('SCBH').AsString;
            with tmpQry do
            begin
              Active:=false;
              SQL.Clear;
              SQL.Add('select isnull(KCRKScan_RFSS.RemQty,KCRKScan_RFSS.Qty) RemQty, KCRKScan_RFSS.CLBH from KCRKScan_RFSS ');
              SQL.Add('Where KCRKScan_RFSS.barcode='''+barcode+''' ');
              active:=true;
              remqty:=tmpQry.fieldbyname('RemQty').AsString;
              Active:=false;
              SQL.Clear;
              SQL.Add('select (max(KCLLS.Qty)-sum(isnull(KCLLScan_RFSS.SCANQty,0))) as SCANQty from KCLLS ');
              SQL.Add('left join KCLLScan_RFSS on KCLLScan_RFSS.LLNO=KCLLS.LLNO and KCLLScan_RFSS.CLBH=KCLLS.CLBH ');
              SQL.Add('						and KCLLScan_RFSS.DFL=KCLLS.DFL and KCLLScan_RFSS.SCBH=KCLLS.SCBH ');
              SQL.Add('Where KCLLS.LLNO='''+KCLLSQry.fieldbyname('LLNO').AsString+''' and KCLLS.CLBH='''+CLBH1+''' ');
              SQL.Add('      and KCLLS.DFL='''+DFL+''' and KCLLS.SCBH='''+SCBH+''' ');
              active:=true;
              SCANQty:=tmpQry.fieldbyname('SCANQty').AsString;
              qty:=FloatToStr(StrToFloat(remqty)-StrToFloat(SCANQty));
              if (StrToFloat(SCANQty)>0) and (StrToFloat(remqty)>0) then
              begin
                with tmpQry do
                begin
                  Active:=false;
                  SQL.Clear;
                  sql.Add('Update KCLLS set SCANUSERID='''+main.Edit1.Text+''', SCANUSERDATE=getdate() ');
                  SQL.Add('Where KCLLS.LLNO='''+KCLLSQry.fieldbyname('LLNO').AsString+''' and KCLLS.CLBH='''+CLBH1+''' ');
                  SQL.Add('      and KCLLS.DFL='''+DFL+''' and KCLLS.SCBH='''+SCBH+''' ');
                  ExecSQL;
                end;
                if (StrToFloat(qty)>0) then
                begin
                  with tmpQry do
                  begin
                    Active:=false;
                    SQL.Clear;
                    sql.Add('Insert into KCLLScan_RFSS (LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE) ');
                    SQL.Add('Values( '''+KCLLSQry.fieldbyname('LLNO').AsString+''','''+CLBH1+''','''+DFL+''','''+SCBH+''', ');
                    SQL.Add('        '''+barcode+''','''+SCANQty+''','''+main.Edit1.Text+''',Getdate()) ');
                    ExecSQL;
                    Active:=false;
                    SQL.Clear;
                    sql.Add('Update KCRKScan_RFSS set RemQty='''+qty+''', CFMID1='''+main.Edit1.Text+''', CFMDATE1=getdate() ');
                    SQL.Add('Where KCRKScan_RFSS.barcode='''+barcode+''' and CFMDATE is not null ');
                    ExecSQL;
                  end;
                end else begin
                  with tmpQry do
                  begin
                    Active:=false;
                    SQL.Clear;
                    sql.Add('Insert into KCLLScan_RFSS (LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE) ');
                    SQL.Add('Values( '''+KCLLSQry.fieldbyname('LLNO').AsString+''','''+CLBH1+''','''+DFL+''','''+SCBH+''', ');
                    SQL.Add('        '''+barcode+''','''+remqty+''','''+main.Edit1.Text+''',Getdate()) ');
                    ExecSQL;
                    Active:=false;
                    SQL.Clear;
                    sql.Add('Update KCRKScan_RFSS set RemQty=0, CFMID1='''+main.Edit1.Text+''', CFMDATE1=getdate() ');
                    SQL.Add('Where KCRKScan_RFSS.barcode='''+barcode+''' and CFMDATE is not null ');
                    ExecSQL;
                  end;
                end;
              end;
            end;
          end;
          KCLLSQry.Next;
        end;
      except
        On E:Exception do
        begin
          Messagedlg('Have wrong, can not save data!'+E.Message,mtwarning,[mbyes],0);
        end;
      end;
    end;
    edit1.Text:='0';
    with Query1 do
    begin
      active:=false;
      SQL.Clear;
      SQL.Add('select KCRKScan_RFSS.SCNO,KCRKScan_RFSS.CLBH,CLZL.YWPM,ZSZL.ZSYWJC,KCRKScan_RFSS.Memo_RY,KCRKScan_RFSS.Memo_Article  ');
      SQL.Add(',KCRKScan_RFSS.Qty,CONVERT(varchar, KCRKScan_RFSS.Pack,101)+'' of '' +CONVERT(varchar,TotalPack,101) as PACK1,KCRKScan_RFSS.Pack  ');
      SQL.Add(',KCRKScan_RF.Memo+'' - ''+KCRKScan_RF.DOCNO as Memo,CONVERT(varchar,isnull(KCRKScan_RFSS.RemQty,KCRKScan_RFSS.Qty),101)+'' ''+CLZL.DWBH DWBH ');
      SQL.Add(',KCRKScan_RFS.CGNO, KCRKScan_RFSS.barcode as QRCode, CONVERT(varchar,month(getdate())) as FIFO, KCRKScan_RFSS.LotNO,KCRKScan_RFS.IDCheck ');
      SQL.Add('from KCRKScan_RFSS    ');
      SQL.Add('left join KCRKScan_RF on KCRKScan_RFSS.SCNO = KCRKScan_RF.SCNO    ');
      SQL.Add('left join KCRKScan_RFS on KCRKScan_RFSS.SCNO = KCRKScan_RFS.SCNO and KCRKScan_RFSS.CLBH = KCRKScan_RFS.CLBH    ');
      SQL.Add('left join (select SCNO,CLBH,count (Pack) as TotalPack   ');
      SQL.Add('  			    from KCRKScan_RFSS   ');
      SQL.Add('  			    group by SCNO,CLBH) KCRKScan  on KCRKScan.SCNO = KCRKScan_RFS.SCNO and KCRKScan.CLBH = KCRKScan_RFS.CLBH   ');
      SQL.Add('left join CGZL on KCRKScan_RFS.CGNO = CGZL.CGNO    ');
      SQL.Add('left join ZSZL on CGZL.ZSBH = ZSZL.ZSDH    ');
      SQL.Add('left join CLZL on KCRKScan_RFSS.CLBH = CLZL.CLDH    ');
      SQL.Add('Where KCRKScan_RFSS.barcode='''+barcode+''' ');
      //funcobj.WriteErrorLog(sql.Text);
      active:=true;
    end;
    QRCodeBitmap1:=TBitmap.Create;
    UpLoadQRImg1(Query1.fieldbyname('QRCode').AsString);
    QRImage2.Picture.Assign(QRCodeBitmap1);
    KCLLSQry.Active:=false;
    KCLLSQry.cachedupdates:=false;
    KCLLSQry.requestlive:=false;
    KCLLSQry.Active:=true;
    Save.Visible:=true;
    LLBarCodeEdit1.Text:='';
    key:=#13;
  end;
end;

procedure TScanQRDelivery.LLBarCodeEditKeyPress(Sender: TObject;
  var Key: Char);
var barcodeLLNO,list: string;
    i:integer;
begin
  if (key='.') or (key=#13) then
  begin
    QRCodeBitmap.Free;
    barcodeLLNO:=LLBarCodeEdit.Text;
    LLNOlist:=LLNOlist+''''+barcodeLLNO+''',';
    list:=Copy(LLNOlist,1,Length(LLNOlist)-1);
    with KCLLSQry do
    begin
      active:=false;
      SQL.Clear;
      SQL.Add('select KCLLS.LLNO,KCLLS.CLBH,KCLLS.SCBH,KCLLS.DFL,CLZL.ywpm,CLZL.DWBH,KCLLS.Qty,SCBLYY.YWSM,DDZL.ARTICLE+''/''+DDZL.BUYNO as SKU ');
      SQL.Add('      ,KCLLS.BLSB,KCLLS.SCANUSERID,KCLLS.SCANUSERDATE, KCLLSAll.Allqty, sum(isnull(KCLLScan_RFSS.SCANQty,0)) LLQty ');
      SQL.Add('	     ,KCLLSAll.Allqty-isnull(KCLLScan_RFSSAll.LLAllqty,0) as TempQty, isnull(KCLLS.SCANQty,0) checkry ');
      SQL.Add('from KCLLS ');
      SQL.Add('left join (select LLNO,CLBH,DFL,SCBH,SCANQty from KCLLScan_RFSS union all ');
      SQL.Add('		        select LLNO,CLBH,DFL,SCBH,SCANQty from KCLLScan_RFSSS) KCLLScan_RFSS ');
      SQL.Add('		       on KCLLScan_RFSS.LLNO=KCLLS.LLNO and KCLLScan_RFSS.CLBH=KCLLS.CLBH ');
      SQL.Add('          and KCLLScan_RFSS.DFL=KCLLS.DFL and KCLLScan_RFSS.SCBH=KCLLS.SCBH ');
      SQL.Add('left join (select CLBH, sum(KCLLS.Qty) Allqty  ');
      SQL.Add('           from KCLLS ');
      SQL.Add('           Where KCLLS.LLNO in ('+list+') ');
      SQL.Add('           group by CLBH) KCLLSAll on KCLLSAll.CLBH=KCLLS.CLBH ');
      SQL.Add('left join (select CLBH, sum(isnull(KCLLScan_RFSS.SCANQty,0)) LLAllqty  ');
      SQL.Add('           from KCLLScan_RFSS ');
      SQL.Add('           Where KCLLScan_RFSS.LLNO in ('+list+') ');
      SQL.Add('           group by CLBH) KCLLScan_RFSSAll on KCLLScan_RFSSAll.CLBH=KCLLS.CLBH ');
      SQL.Add('left join CLZL on CLZL.cldh=KCLLS.CLBH ');
      SQL.Add('left join DDZL on DDZL.DDBH=KCLLS.SCBH ');
      SQL.Add('left join SCBLYY on SCBLYY.YYBH=KCLLS.YYBH ');
      SQL.Add('Where KCLLS.LLNO in ('+list+') ');
      SQL.Add('group by KCLLS.LLNO,KCLLS.CLBH,KCLLS.SCBH,KCLLS.DFL,CLZL.ywpm,CLZL.DWBH,DDZL.ARTICLE,KCLLS.Qty,DDZL.BUYNO ');
      SQL.Add('         ,KCLLS.SCANQty,SCBLYY.YWSM ,KCLLS.BLSB,KCLLS.SCANUSERID,KCLLS.SCANUSERDATE,KCLLSAll.Allqty,KCLLScan_RFSSAll.LLAllqty ');
      SQL.Add('order by KCLLS.CLBH ');
      //funcobj.WriteErrorLog(sql.Text);
      active:=true;
    end;
    QRCodeBitmap:=TBitmap.Create;
    UpLoadQRImg(KCLLSQry.fieldbyname('LLNO').AsString);
    QRImage1.Picture.Assign(QRCodeBitmap);
    LLBarCodeEdit.Text:='';
    key:=#13;
  end;
end;

procedure TScanQRDelivery.imgSave00Click(Sender: TObject);
var writeFile:textfile;
    LabelInfo:string;
begin
  //PrintLabel_LL:=TPrintLabel_LL.Create(nil);
  //PrintLabel_LL.quickrep1.prepare;
  //PrintLabel_LL.quickrep1.preview;
  //PrintLabel_LL.Free;
  if Query1.RecordCount>0 then
  begin
    if Fileexists(extractfilepath(application.ExeName)+'\label_N8A.txt') then
    begin
      assignfile(writeFile,ExtractFilePath(Application.ExeName)+'label_N8A.txt');
      rewrite(writeFile);
      Query1.first;
      while not Query1.eof do
      begin
        LabelInfo:='';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('QRCode').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('CLBH').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('YWPM').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('ZSYWJC').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('CGNO').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('Memo_RY').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('Memo_Article').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('DWBH').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('PACK1').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('Memo').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('FIFO').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('LotNO').AsString)+'|';
        LabelInfo:=LabelInfo+trim(Query1.fieldByName('IDCheck').AsString)+'|';
        append(writeFile);
        writeln(writeFile, LabelInfo);
          //
        Query1.next;
      end;
      closefile(writeFile);
    end;
    //
    if fileexists(extractfilepath(application.ExeName)+'Label_N8A.btw') then
    begin
      ShellExecute(Handle,'open',PChar(ExtractFilePath(Application.ExeName)+'Label_N8A.btw'),nil,pchar(ExtractFilePath(Application.ExeName)),SW_SHOW);
      //close;
    end
    else  begin
      Showmessage('Pls setup the program first. Abort');
    end;
  end else begin
    CopyFile(PChar(ScanQR.WHFilePath+'label_N8A.txt'),PChar(extractfilepath(application.ExeName)+'Label_N8A.txt'), false);
    if (MessageDlg('Xoa luu du lieu nhan in?', mtConfirmation, [mbYes, mbNo], 0) = mrYes) then
    begin
      if Fileexists(ScanQR.WHFilePath+'label_N8A.txt') then
      begin
        assignfile(writeFile,ScanQR.WHFilePath+'label_N8A.txt');
        rewrite(writeFile);
        closefile(writeFile);
      end;
    end;
    if fileexists(extractfilepath(application.ExeName)+'Label_N8A.btw') then
    begin
      ShellExecute(Handle,'open',PChar(ExtractFilePath(Application.ExeName)+'Label_N8A.btw'),nil,pchar(ExtractFilePath(Application.ExeName)),SW_SHOW);
      //close;
    end
    else  begin
      Showmessage('Pls setup the program first. Abort');
    end;
  end;
  //
end;

procedure TScanQRDelivery.KCLLSQryAfterOpen(DataSet: TDataSet);
begin
  if KCLLSQry.RecordCount>0 then
    edit1.ReadOnly:=true
  else edit1.ReadOnly:=false;
end;

procedure TScanQRDelivery.TrackClick(Sender: TObject);
begin
  ScanMaterialQty:=TScanMaterialQty.Create(self);
  ScanMaterialQty.Show;
end;

procedure TScanQRDelivery.SaveClick(Sender: TObject);
var writeFile:textfile;
    LabelInfo:string;
begin
  if Fileexists(ScanQR.WHFilePath+'label_N8A.txt') then
  begin
    assignfile(writeFile,ScanQR.WHFilePath+'label_N8A.txt');
    append(writeFile);
    Query1.first;
    while not Query1.eof do
    begin
      LabelInfo:='';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('QRCode').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('CLBH').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('YWPM').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('ZSYWJC').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('CGNO').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('Memo_RY').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('Memo_Article').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('DWBH').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('PACK1').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('Memo').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('FIFO').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('LotNO').AsString)+'|';
      LabelInfo:=LabelInfo+trim(Query1.fieldByName('IDCheck').AsString)+'|';
      append(writeFile);
      writeln(writeFile, LabelInfo);
        //
      Query1.next;
    end;
    closefile(writeFile);
    Showmessage('Da luu nhan thanh cong !');
  end;
end;

procedure TScanQRDelivery.DBGridEh2GetCellParams(Sender: TObject;
  Column: TColumnEh; AFont: TFont; var Background: TColor;
  State: TGridDrawState);
begin
  if KCLLSQry.FieldByName('checkry').Value=1 then
  begin
    DBGridEh2.canvas.font.color:=clred;
  end;
end;

procedure TScanQRDelivery.DBGridEh2CellClick(Column: TColumnEh);
begin
  if KCLLSQry.fieldbyname('checkry').Value=1 then
  begin
    with KCLLSQry do
    begin
      requestlive:=true;
      cachedupdates:=true;
      edit;
      fieldbyname('checkry').Value:=0;
      UpdateSQL1.apply(ukmodify);
    end;
  end else
  begin
    with KCLLSQry do
    begin
      requestlive:=true;
      cachedupdates:=true;
      edit;
      fieldbyname('checkry').Value:=1;
      UpdateSQL1.apply(ukmodify);
    end;
  end;
end;

procedure TScanQRDelivery.ReportDblClick(Sender: TObject);
begin
  if (MessageDlg('Reset scan phieu xuat?', mtConfirmation, [mbYes, mbNo], 0) = mrYes) then
  begin
    LLBarCodeEdit1.Visible:=false;
    LLBarCodeEdit.Visible:=true;
    LLBarCodeEdit.SetFocus;
    edit1.ReadOnly:=false;
    LLNOlist :='';
    with KCLLSQry do
    begin
      active:=false;
      SQL.Clear;
      SQL.Add('select KCLLS.LLNO,KCLLS.CLBH,KCLLS.SCBH,KCLLS.DFL,CLZL.ywpm,CLZL.DWBH,KCLLS.Qty,SCBLYY.YWSM,DDZL.ARTICLE+''/''+DDZL.BUYNO as SKU ');
      SQL.Add('      ,KCLLS.BLSB,KCLLS.SCANUSERID,KCLLS.SCANUSERDATE, KCLLSAll.Allqty, sum(isnull(KCLLScan_RFSS.SCANQty,0)) LLQty ');
      SQL.Add('	     ,KCLLSAll.Allqty-isnull(KCLLScan_RFSSAll.LLAllqty,0) as TempQty, isnull(KCLLS.SCANQty,0) checkry ');
      SQL.Add('from KCLLS ');
      SQL.Add('left join (select LLNO,CLBH,DFL,SCBH,SCANQty from KCLLScan_RFSS union all ');
      SQL.Add('		        select LLNO,CLBH,DFL,SCBH,SCANQty from KCLLScan_RFSSS) KCLLScan_RFSS ');
      SQL.Add('		       on KCLLScan_RFSS.LLNO=KCLLS.LLNO and KCLLScan_RFSS.CLBH=KCLLS.CLBH ');
      SQL.Add('          and KCLLScan_RFSS.DFL=KCLLS.DFL and KCLLScan_RFSS.SCBH=KCLLS.SCBH ');
      SQL.Add('left join (select CLBH, sum(KCLLS.Qty) Allqty  ');
      SQL.Add('           from KCLLS ');
      SQL.Add('           Where 1=2 ');
      SQL.Add('           group by CLBH) KCLLSAll on KCLLSAll.CLBH=KCLLS.CLBH ');
      SQL.Add('left join (select CLBH, sum(isnull(KCLLScan_RFSS.SCANQty,0)) LLAllqty  ');
      SQL.Add('           from KCLLScan_RFSS ');
      SQL.Add('           Where 1=2 ');
      SQL.Add('           group by CLBH) KCLLScan_RFSSAll on KCLLScan_RFSSAll.CLBH=KCLLS.CLBH ');
      SQL.Add('left join CLZL on CLZL.cldh=KCLLS.CLBH ');
      SQL.Add('left join DDZL on DDZL.DDBH=KCLLS.SCBH ');
      SQL.Add('left join SCBLYY on SCBLYY.YYBH=KCLLS.YYBH ');
      SQL.Add('Where 1=2 ');
      SQL.Add('group by KCLLS.LLNO,KCLLS.CLBH,KCLLS.SCBH,KCLLS.DFL,CLZL.ywpm,CLZL.DWBH,DDZL.ARTICLE,KCLLS.Qty,DDZL.BUYNO ');
      SQL.Add('         ,KCLLS.SCANQty,SCBLYY.YWSM ,KCLLS.BLSB,KCLLS.SCANUSERID,KCLLS.SCANUSERDATE,KCLLSAll.Allqty,KCLLScan_RFSSAll.LLAllqty ');
      SQL.Add('order by KCLLS.CLBH ');
      //funcobj.WriteErrorLog(sql.Text);
      active:=true;
    end;
  end;
end;

procedure TScanQRDelivery.DBGridEh2DblClick(Sender: TObject);
begin
  if (KCLLSQry.fieldbyname('LLNO').AsString<>'') and (KCLLSQry.fieldbyname('Qty').Value<=KCLLSQry.fieldbyname('LLQty').Value) then
  begin
    ScanQRDeliveryOWE:=TScanQRDeliveryOWE.Create(self);

    ScanQRDeliveryOWE.Height:=Round(Screen.Height*3/4);
    ScanQRDeliveryOWE.Width:=Round(Screen.Width*3/4);
    ScanQRDeliveryOWE.Left:=Round((Screen.Width-ScanQRDeliveryOWE.Width)/2);
    ScanQRDeliveryOWE.Top:=Round((Screen.Height-ScanQRDeliveryOWE.Height)/2);

    ScanQRDeliveryOWE.imgClose00.Width:=Round(Screen.Width/28);
    ScanQRDeliveryOWE.imgClose00.Height:=Round(Screen.Width/28);
    ScanQRDeliveryOWE.imgClose00.Left:=Round(ScanQRDeliveryOWE.Width-ScanQRDeliveryOWE.imgClose00.Width);
    ScanQRDeliveryOWE.imgClose00.Top:=Round(Screen.DesktopTop+5);

    ScanQRDeliveryOWE.LLNO:=KCLLSQry.fieldbyname('LLNO').AsString;
    ScanQRDeliveryOWE.CLBH:=KCLLSQry.fieldbyname('CLBH').AsString;
    ScanQRDeliveryOWE.DFL:=KCLLSQry.fieldbyname('DFL').AsString;
    ScanQRDeliveryOWE.SCBH:=KCLLSQry.fieldbyname('SCBH').AsString;

    ScanQRDeliveryOWE.LLBarCodeEdit.SetFocus;

    ScanQRDeliveryOWE.Show;
  end;
end;

procedure TScanQRDelivery.ImageOUTClick(Sender: TObject);
begin
  ScanQRDeliveryPAY:=TScanQRDeliveryPAY.Create(self);
  ScanQRDeliveryPAY.Show;
end;

end.
