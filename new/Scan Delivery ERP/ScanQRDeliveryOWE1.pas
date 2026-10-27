unit ScanQRDeliveryOWE1;

interface

uses
  Windows, Messages, SysUtils, Variants, Classes, Graphics, Controls, Forms,
  Dialogs, DBTables, DB, ScktComp, Menus, GridsEh, DBGridEh, ExtCtrls,pngimage,
  StdCtrls, Mask, DBCtrls, Buttons, ComCtrls, comobj, jpeg, Spin, DelphiZXIngQRCode,
  QRCtrls, QuickRpt,ShellApi, Printers, RpCon, RpConDS, RpConBDE, RpDefine, RpRave,
  GIFImage, dateutils;

type
  TScanQRDeliveryOWE = class(TForm)
    Panel3: TPanel;
    Image2: TImage;
    Label2: TLabel;
    imgClose00: TImage;
    Edit1: TEdit;
    Save: TImage;
    DBGridEh2: TDBGridEh;
    KCLLSQry: TQuery;
    KCLLSQryCLBH: TStringField;
    KCLLSQryLLNO: TStringField;
    KCLLSQrySCBH: TStringField;
    KCLLSQryDFL: TStringField;
    KCLLSQrySCANUSERID: TStringField;
    KCLLSQrySCANUSERDATE: TDateTimeField;
    DS2: TDataSource;
    LLBarCodeEdit: TEdit;
    tmpQry: TQuery;
    KCLLSQrybarcode: TStringField;
    KCLLSQrySCANQty: TCurrencyField;
    KCLLSQryLLNO_OWE: TStringField;
    KCLLSQryywpm: TStringField;
    KCLLSQrydwbh: TStringField;
    KCLLSQrySKU: TStringField;
    UpdKCLLS: TUpdateSQL;
    BDelRec: TQuery;
    procedure FormDestroy(Sender: TObject);
    procedure FormClose(Sender: TObject; var Action: TCloseAction);
    procedure FormCreate(Sender: TObject);
    procedure imgClose00Click(Sender: TObject);
    procedure LLBarCodeEditKeyPress(Sender: TObject; var Key: Char);
    procedure SaveClick(Sender: TObject);
    procedure Edit1KeyPress(Sender: TObject; var Key: Char);
  private
    { Private declarations }
  public
    LLNO, CLBH, DFL, SCBH: string;
    { Public declarations }
  end;

var
  ScanQRDeliveryOWE: TScanQRDeliveryOWE;

implementation

uses main1, ScanQRDelivery1, FunUnit;

{$R *.dfm}

procedure TScanQRDeliveryOWE.FormDestroy(Sender: TObject);
begin
  ScanQRDeliveryOWE:=nil;
end;

procedure TScanQRDeliveryOWE.FormClose(Sender: TObject;
  var Action: TCloseAction);
begin
  Action:=Cafree;
end;

procedure TScanQRDeliveryOWE.FormCreate(Sender: TObject);
var imout, imsave: TPNGObject;
begin
  imout := TPNGObject.Create;
  imout.LoadFromFile('btnlogout.png');
  imgClose00.Picture.Assign(imout);

  imsave := TPNGObject.Create;
  imsave.LoadFromFile('btnsave.png');
  Save.Picture.Assign(imsave);
  
  imout.free;
  imsave.free;

  KCLLSQry.Active:=true;
end;

procedure TScanQRDeliveryOWE.imgClose00Click(Sender: TObject);
begin
  ScanQRDelivery.KCLLSQry.Active:=false;
  ScanQRDelivery.KCLLSQry.Active:=true;
  Close;
end;

procedure TScanQRDeliveryOWE.LLBarCodeEditKeyPress(Sender: TObject;
  var Key: Char);
var barcode, remqty, CFMDATE, qty, y,m, a: string;
begin
  if (key='.') or (key=#13) then
  begin
    if edit1.Text='' then edit1.Text:='0';
    ScanQRDelivery.QRCodeBitmap1.Free;
    barcode:=LLBarCodeEdit.Text;
    with tmpQry do
    begin
      Active:=false;
      SQL.Clear;
      SQL.Add('select isnull(KCRKScan_RFSS.RemQty,KCRKScan_RFSS.Qty) RemQty, CFMDATE from KCRKScan_RFSS ');
      SQL.Add('Where KCRKScan_RFSS.barcode='''+barcode+''' ');
      active:=true;
      remqty:=tmpQry.fieldbyname('RemQty').AsString;
      CFMDATE:=tmpQry.fieldbyname('CFMDATE').AsString;
    end;
    if CFMDATE='' then
    begin
      Showmessage('Vui long scan nhap kho !!!');
      edit1.Text:='0';
      KCLLSQry.Active:=false;
      KCLLSQry.cachedupdates:=false;
      KCLLSQry.requestlive:=false;
      KCLLSQry.Active:=true;
      LLBarCodeEdit.Text:='';
      key:=#13;
      Exit;
    end;
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
          SQL.Add('Values( '''+a+''','''+LLNO+''','''+CLBH+''','''+DFL+''','''+SCBH+''', ');
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
          SQL.Add('Values( '''+a+''','''+LLNO+''','''+CLBH+''','''+DFL+''','''+SCBH+''', ');
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
    KCLLSQry.Active:=false;
    KCLLSQry.Active:=true;
    with ScanQRDelivery.Query1 do
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
    ScanQRDelivery.QRCodeBitmap1:=TBitmap.Create;
    ScanQRDelivery.UpLoadQRImg1(barcode);
    ScanQRDelivery.QRImage2.Picture.Assign(ScanQRDelivery.QRCodeBitmap1);
    ScanQRDelivery.imgSave00.Visible:=true;
    ScanQRDelivery.Save.Visible:=true;
    edit1.Text:='0';
    LLBarCodeEdit.Text:='';
    key:=#13;
  end;
end;

procedure TScanQRDeliveryOWE.SaveClick(Sender: TObject);
var i:integer;
begin
  try
    KCLLSQry.first;
    for i:=1 to KCLLSQry.RecordCount do
    begin
      if (KCLLSQry.fieldbyname('SCANQty').AsString='0') then
      begin
        with BDelRec do
        begin
          active:=false;
          sql.Clear;
          sql.add('insert into BDelRec ');
          sql.add('(TableName ,TNO,CLBH,OldID,DelID,DelDate)');
          sql.add('values (''KCLLScan_RFSSS'','''+KCLLSQry.fieldbyname('LLNO_OWE').Value+'*'+KCLLSQry.fieldbyname('LLNO').Value+'*'+KCLLSQry.fieldbyname('SCBH').Value+'''');
          sql.add(','''+KCLLSQry.fieldbyname('CLBH').Value+'*'+KCLLSQry.fieldbyname('barcode').Value+'*'+FloatToStr(KCLLSQry.fieldbyname('SCANQty').OldValue)+''',');
          sql.add(''''+KCLLSQry.fieldbyname('SCANUSERID').Value+''','''+main.Edit1.Text+''',getdate())');
          ExecSQL;
          active:=false;
        end;
        UpdKCLLS.apply(ukdelete);
        with tmpQry do
        begin          Active:=false;          SQL.Clear;          SQL.Add('Update KCRKScan_RFSS set RemQty=RemQty+'+FloatToStr(KCLLSQry.fieldbyname('SCANQty').OldValue)+' ');          SQL.Add('Where KCRKScan_RFSS.barcode='''+KCLLSQry.FieldByName('barcode').AsString+''' and CFMDATE is not null ');          ExecSQL;        end;
      end;
      KCLLSQry.Next;
    end;

    Showmessage('OK');
    KCLLSQry.Active:=false;
    KCLLSQry.Active:=true;
  except
    On E:Exception do
    begin
      Messagedlg('Have wrong, can not save data!'+E.Message,mtwarning,[mbyes],0);
    end;
  end;
end;

procedure TScanQRDeliveryOWE.Edit1KeyPress(Sender: TObject; var Key: Char);
  var s:string;
begin
  s := ('1234567890.'#8); //Add chars you want to allow
  if pos(key,s) =0 then Key:=#0;
end;

end.
