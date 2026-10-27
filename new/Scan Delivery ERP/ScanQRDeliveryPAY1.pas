unit ScanQRDeliveryPAY1;

interface

uses
  Windows, Messages, SysUtils, Variants, Classes, Graphics, Controls, Forms,
  Dialogs, DBTables, DB, ScktComp, Menus, GridsEh, DBGridEh, ExtCtrls,pngimage,
  StdCtrls, Mask, DBCtrls, Buttons, ComCtrls, comobj, jpeg, Spin, DelphiZXIngQRCode,
  QRCtrls, QuickRpt,ShellApi, Printers, RpCon, RpConDS, RpConBDE, RpDefine, RpRave,
  GIFImage;

type
  TScanQRDeliveryPAY = class(TForm)
    Panel1: TPanel;
    Image1: TImage;
    imgClose00: TImage;
    Label6: TLabel;
    Label2: TLabel;
    Label5: TLabel;
    Label7: TLabel;
    Report: TImage;
    SCNOEdit: TEdit;
    DTP1: TDateTimePicker;
    DTP2: TDateTimePicker;
    CLEdit: TEdit;
    CLNOEdit: TEdit;
    DBGridEh2: TDBGridEh;
    KCLLScan_RFSSS: TQuery;
    KCLLScan_RFSSSLLNO_OWE: TStringField;
    KCLLScan_RFSSSLLNO: TStringField;
    KCLLScan_RFSSSCLBH: TStringField;
    KCLLScan_RFSSSywpm: TStringField;
    KCLLScan_RFSSSdwbh: TStringField;
    KCLLScan_RFSSSDFL: TStringField;
    KCLLScan_RFSSSSCBH: TStringField;
    KCLLScan_RFSSSSKU: TStringField;
    KCLLScan_RFSSSSCANQty: TCurrencyField;
    KCLLScan_RFSSSbarcode: TStringField;
    KCLLScan_RFSSSSCANUSERID: TStringField;
    KCLLScan_RFSSSSCANUSERDATE: TDateTimeField;
    UpdKCLLSS: TUpdateSQL;
    DS2: TDataSource;
    BDelRec: TQuery;
    Panel2: TPanel;
    Panel3: TPanel;
    Panel4: TPanel;
    Image2: TImage;
    Panel5: TPanel;
    Panel7: TPanel;
    Image3: TImage;
    Save: TImage;
    LLBarCodeEdit: TEdit;
    KCLLScan_RFSS: TQuery;
    KCLLScan_RFSSLLNO: TStringField;
    KCLLScan_RFSSCLBH: TStringField;
    KCLLScan_RFSSDFL: TStringField;
    KCLLScan_RFSSSCBH: TStringField;
    KCLLScan_RFSSbarcode: TStringField;
    KCLLScan_RFSSSCANQty: TCurrencyField;
    KCLLScan_RFSSSCANUSERID: TStringField;
    KCLLScan_RFSSSCANUSERDATE: TDateTimeField;
    UpdKCLLS: TUpdateSQL;
    DS3: TDataSource;
    DBGridEh1: TDBGridEh;
    DS1: TDataSource;
    DBGridEh3: TDBGridEh;
    KCLLSQry: TQuery;
    KCLLSQryCLBH: TStringField;
    KCLLSQryywpm: TStringField;
    KCLLSQryDWBH: TStringField;
    KCLLSQryTempQty: TCurrencyField;
    KCLLSQryQty: TCurrencyField;
    KCLLSQryAllqty: TCurrencyField;
    KCLLSQryLLQty: TCurrencyField;
    KCLLSQrySKU: TStringField;
    KCLLSQryYWSM: TStringField;
    KCLLSQryBLSB: TStringField;
    KCLLSQryLLNO: TStringField;
    KCLLSQrySCBH: TStringField;
    KCLLSQryDFL: TStringField;
    KCLLSQrySCANUSERID: TStringField;
    KCLLSQrySCANUSERDATE: TDateTimeField;
    Label1: TLabel;
    UpdateSQL1: TUpdateSQL;
    tmpQry: TQuery;
    KCLLScan_RFSSLLNO_OWE: TStringField;
    KCLLSQrycheckry: TFloatField;
    procedure FormDestroy(Sender: TObject);
    procedure FormClose(Sender: TObject; var Action: TCloseAction);
    procedure FormCreate(Sender: TObject);
    procedure imgClose00Click(Sender: TObject);
    procedure ReportClick(Sender: TObject);
    procedure LLBarCodeEditKeyPress(Sender: TObject; var Key: Char);
    procedure DBGridEh3GetCellParams(Sender: TObject; Column: TColumnEh;
      AFont: TFont; var Background: TColor; State: TGridDrawState);
    procedure DBGridEh2DblClick(Sender: TObject);
    procedure DBGridEh3DblClick(Sender: TObject);
    procedure SaveClick(Sender: TObject);
  private
    LLNOlist:string;
    { Private declarations }
  public
    { Public declarations }
  end;

var
  ScanQRDeliveryPAY: TScanQRDeliveryPAY;

implementation

{$R *.dfm}

uses main1, ScanQRDelivery1, FunUnit;

procedure TScanQRDeliveryPAY.FormDestroy(Sender: TObject);
begin
  ScanQRDeliveryPAY:=nil;
end;

procedure TScanQRDeliveryPAY.FormClose(Sender: TObject;
  var Action: TCloseAction);
begin
  Action:=Cafree;
end;

procedure TScanQRDeliveryPAY.FormCreate(Sender: TObject);
var imout, imquery, imsave: TPNGObject;
begin
  imout := TPNGObject.Create;
  imout.LoadFromFile('btnlogout.png');
  imgClose00.Picture.Assign(imout);

  imquery := TPNGObject.Create;
  imquery.LoadFromFile('btnquery.png');
  Report.Picture.Assign(imquery);

  imsave := TPNGObject.Create;
  imsave.LoadFromFile('btnsave.png');
  Save.Picture.Assign(imsave);

  DBGridEh2.Width:=Round(Screen.Width/2);

  imgClose00.Width:=Round(Screen.Width/28);
  imgClose00.Height:=Round(Screen.Width/28);
  imgClose00.Left:=Round(Screen.Width-imgClose00.Width);
  imgClose00.Top:=Round(Screen.DesktopTop+5);

  Panel3.Height:=Round((Screen.Height-Panel1.Height)/2);

  DTP1.Date:=Date()-3;
  DTP2.date:=Date();

  imout.free;
  imquery.free;
  imsave.free;
end;

procedure TScanQRDeliveryPAY.imgClose00Click(Sender: TObject);
begin
  if MessageDlg('Ban co muon dong khong?',mtConfirmation, [mbYes, mbNo], 0) = mrYes then
  begin
    Close;
  end;
end;

procedure TScanQRDeliveryPAY.ReportClick(Sender: TObject);
begin
  with KCLLScan_RFSSS do
  begin
   Active:=false;
   SQL.Clear;
   SQL.Add('select KCLLScan_RFSSS.*, ywpm, dwbh,DDZL.ARTICLE+''/''+DDZL.BUYNO as SKU ');
   SQL.Add('from KCLLScan_RFSSS ');
   SQL.Add('left join CLZL on CLZL.cldh=KCLLScan_RFSSS.CLBH ');
   SQL.Add('left join DDZL on DDZL.DDBH=KCLLScan_RFSSS.SCBH ');
   SQL.Add('Where KCLLScan_RFSSS.SCANQty<>0 and convert(smalldatetime,convert(varchar,KCLLScan_RFSSS.SCANUSERDATE,111)) between ');
   SQL.add('      '''+formatdatetime('yyyy/MM/dd',DTP1.Date) +''''+' and '+''''+formatdatetime('yyyy/MM/dd',DTP2.Date) +'''');
   if SCNOEdit.Text<>'' then
    SQL.Add('   and KCLLScan_RFSSS.SCBH like '''+SCNOEdit.Text+'%'' ');
   if CLNOEdit.Text<>'' then
     SQl.Add('   and KCLLScan_RFSSS.CLBH like '''+CLNOEdit.Text+'%'' ');
   if CLEdit.Text<>'' then
     SQl.Add('   and CLZL.YWPM like ''%'+CLEdit.Text+'%'' ');
   SQL.Add('order by KCLLScan_RFSSS.LLNO_OWE desc, KCLLScan_RFSSS.CLBH desc ');
   //funcobj.WriteErrorLog(sql.Text);
   active:=true;
  end;
end;

procedure TScanQRDeliveryPAY.LLBarCodeEditKeyPress(Sender: TObject;
  var Key: Char);
var barcodeLLNO,list: string;
begin
  if (key='.') or (key=#13) then
  begin
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
    LLBarCodeEdit.Text:='';
    KCLLScan_RFSS.Active:=false;
    KCLLScan_RFSS.Active:=true;
    key:=#13;
  end;
end;

procedure TScanQRDeliveryPAY.DBGridEh3GetCellParams(Sender: TObject;
  Column: TColumnEh; AFont: TFont; var Background: TColor;
  State: TGridDrawState);
begin
  if KCLLSQry.FieldByName('checkry').Value=1 then
  begin
    DBGridEh3.canvas.font.color:=clred;
  end;
end;

procedure TScanQRDeliveryPAY.DBGridEh2DblClick(Sender: TObject);
var LLNO_OWE, CLBH, DFL, SCBH, barcode, SCANQty, SCANQty1, qty: string;
    i: integer;
begin
  if (KCLLScan_RFSSS.fieldbyname('LLNO_OWE').AsString<>'') and (KCLLSQry.fieldbyname('CLBH').AsString<>'') then
  begin
    with KCLLScan_RFSSS do
    begin
      LLNO_OWE:=fieldbyname('LLNO_OWE').AsString;
      CLBH:=fieldbyname('CLBH').AsString;
      barcode:=fieldbyname('barcode').AsString;
      SCANQty:=fieldbyname('SCANQty').AsString;
    end;
    KCLLSQry.first;
    for i:=1 to KCLLSQry.RecordCount do
    begin
      if (CLBH=KCLLSQry.fieldbyname('CLBH').AsString) and (KCLLSQry.fieldbyname('checkry').Value=1)
          and (KCLLSQry.fieldbyname('LLQty').Value<KCLLSQry.fieldbyname('Qty').Value) then
      begin
        DFL:=KCLLSQry.fieldbyname('DFL').AsString;
        SCBH:=KCLLSQry.fieldbyname('SCBH').AsString;
        SCANQty1:=FloatToStr(StrToFloat(KCLLSQry.fieldbyname('Qty').Value)-StrToFloat(KCLLSQry.fieldbyname('LLQty').Value));
        qty:=FloatToStr(StrToFloat(SCANQty)-StrToFloat(SCANQty1));
        if (StrToFloat(SCANQty1)>0) and (StrToFloat(SCANQty)>0) then
        begin
          with tmpQry do
          begin
            Active:=false;
            SQL.Clear;
            sql.Add('Update KCLLS set SCANUSERID='''+main.Edit1.Text+''', SCANUSERDATE=getdate() ');
            SQL.Add('Where KCLLS.LLNO='''+KCLLSQry.fieldbyname('LLNO').AsString+''' and KCLLS.CLBH='''+CLBH+''' ');
            SQL.Add('      and KCLLS.DFL='''+DFL+''' and KCLLS.SCBH='''+SCBH+''' ');
            ExecSQL;
          end;
          if (StrToFloat(qty)>0) then
          begin
            with tmpQry do
            begin
              Active:=false;
              SQL.Clear;
              sql.Add('Insert into KCLLScan_RFSS (LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE, LLNO_OWE) ');
              SQL.Add('Values( '''+KCLLSQry.fieldbyname('LLNO').AsString+''','''+CLBH+''','''+DFL+''','''+SCBH+''', ');
              SQL.Add('        '''+barcode+''','''+SCANQty1+''','''+main.Edit1.Text+''',Getdate(),'''+LLNO_OWE+''') ');
              ExecSQL;
              Active:=false;
              SQL.Clear;
              sql.Add('Update KCLLScan_RFSSS set SCANQty='''+qty+''' ');
              SQL.Add('Where KCLLScan_RFSSS.LLNO_OWE='''+LLNO_OWE+''' ');
              ExecSQL;
              SCANQty:=qty;
            end;
          end else begin
            with tmpQry do
            begin
              Active:=false;
              SQL.Clear;
              sql.Add('Insert into KCLLScan_RFSS (LLNO, CLBH, DFL, SCBH, barcode, SCANQty, SCANUSERID, SCANUSERDATE, LLNO_OWE) ');
              SQL.Add('Values( '''+KCLLSQry.fieldbyname('LLNO').AsString+''','''+CLBH+''','''+DFL+''','''+SCBH+''', ');
              SQL.Add('        '''+barcode+''','''+SCANQty+''','''+main.Edit1.Text+''',Getdate(),'''+LLNO_OWE+''') ');
              ExecSQL;
              Active:=false;
              SQL.Clear;
              sql.Add('Update KCLLScan_RFSSS set SCANQty=0 ');
              SQL.Add('Where KCLLScan_RFSSS.LLNO_OWE='''+LLNO_OWE+''' ');
              ExecSQL;
            end;
          end;
        end;
      end;
      KCLLSQry.Next;
    end;
    KCLLScan_RFSSS.Active:=false;
    KCLLScan_RFSSS.Active:=true;
    KCLLSQry.Active:=false;
    KCLLSQry.Active:=true;
    KCLLScan_RFSS.Active:=false;
    KCLLScan_RFSS.Active:=true;
  end;
end;

procedure TScanQRDeliveryPAY.DBGridEh3DblClick(Sender: TObject);
begin
  if LLNOlist ='' then
  begin
    KCLLSQry.active:=true;
  end;
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

procedure TScanQRDeliveryPAY.SaveClick(Sender: TObject);
var i:integer;
    qty:string;
begin
  try
    KCLLScan_RFSS.first;
    for i:=1 to KCLLScan_RFSS.RecordCount do
    begin
      qty:='0';
      if KCLLScan_RFSS.fieldbyname('SCANQty').AsString='0'then
      begin
        with BDelRec do
        begin
          active:=false;
          sql.Clear;
          sql.add('insert into BDelRec ');
          sql.add('(TableName ,TNO,CLBH,OldID,DelID,DelDate)');
          sql.add('values (''KCLLScan_RFSS'','''+KCLLScan_RFSS.fieldbyname('LLNO').Value+'*'+KCLLScan_RFSS.fieldbyname('SCBH').Value+'''');
          sql.add(','''+KCLLScan_RFSS.fieldbyname('CLBH').Value+'*'+KCLLScan_RFSS.fieldbyname('barcode').Value+'*'+FloatToStr(KCLLScan_RFSS.fieldbyname('SCANQty').OldValue)+''',');
          sql.add(''''+KCLLScan_RFSS.fieldbyname('SCANUSERID').Value+''','''+main.Edit1.Text+''',getdate())');
          ExecSQL;
          active:=false;
        end;
        with tmpQry do
        begin          Active:=false;          SQL.Clear;          sql.Add('Update KCLLScan_RFSSS set SCANQty=SCANQty+'+FloatToStr(KCLLScan_RFSS.fieldbyname('SCANQty').OldValue)+' ');          SQL.Add('Where KCLLScan_RFSSS.LLNO_OWE='''+KCLLScan_RFSS.fieldbyname('LLNO_OWE').Value+''' ');          ExecSQL;        end;
        UpdKCLLS.apply(ukdelete);
      end else
      begin
        KCLLScan_RFSS.edit;
        KCLLScan_RFSS.FieldByName('SCANUSERID').Value:=main.edit1.text;
        KCLLScan_RFSS.FieldByName('SCANUSERDATE').Value:=date;
        qty:=FloatToStr(KCLLScan_RFSS.fieldbyname('SCANQty').OldValue-KCLLScan_RFSS.fieldbyname('SCANQty').Value);
        with tmpQry do
        begin          Active:=false;          SQL.Clear;          sql.Add('Update KCLLScan_RFSSS set SCANQty=SCANQty+'+qty+' ');          SQL.Add('Where KCLLScan_RFSSS.LLNO_OWE='''+KCLLScan_RFSS.fieldbyname('LLNO_OWE').Value+''' ');          ExecSQL;        end;
        UpdKCLLS.apply(ukmodify);
      end;
      KCLLScan_RFSS.Next;
    end;

    Showmessage('OK');
    KCLLScan_RFSSS.Active:=false;
    KCLLScan_RFSSS.Active:=true;
    KCLLSQry.Active:=false;
    KCLLSQry.Active:=true;
    KCLLScan_RFSS.Active:=false;
    KCLLScan_RFSS.Active:=true;
  except
    On E:Exception do
    begin
      Messagedlg('Have wrong, can not save data!'+E.Message,mtwarning,[mbyes],0);
    end;
  end;
end;

end.
