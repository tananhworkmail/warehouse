SET XACT_ABORT ON;
BEGIN TRANSACTION;

/* Remove only duplicate rows created by the first legacy migration, when the
   old varchar DepNo column converted N'Ủng' to '?ng'. The JSON source remains
   intact and these two forms are inserted again after this repair. */
DECLARE @InvalidReports TABLE (ReportID varchar(50) PRIMARY KEY);
INSERT INTO @InvalidReports (ReportID)
SELECT report.ReportID
FROM dbo.GoreTex_Waterproof_Report report
WHERE report.UpdatedBy = 'legacy-migration'
  AND report.DepNo = '?ng'
  AND EXISTS (
      SELECT 1
      FROM dbo.gore_tex_waterproof_forms legacy
      WHERE legacy.line = N'Ủng'
        AND legacy.style_name = report.StyleName
  );

DELETE snapshot
FROM dbo.gore_tex_dashboard_defect_snapshots snapshot
JOIN @InvalidReports invalid ON invalid.ReportID = snapshot.source_record_key
WHERE snapshot.period_type = 'D';

DELETE detailRevision
FROM dbo.GoreTex_Waterproof_Report_Detail_Revision detailRevision
JOIN dbo.GoreTex_Waterproof_Report_Revision revision ON revision.RevisionID = detailRevision.RevisionID
JOIN @InvalidReports invalid ON invalid.ReportID = revision.ReportID;

DELETE revision
FROM dbo.GoreTex_Waterproof_Report_Revision revision
JOIN @InvalidReports invalid ON invalid.ReportID = revision.ReportID;

DELETE detail
FROM dbo.GoreTex_Waterproof_Report_Detail detail
JOIN @InvalidReports invalid ON invalid.ReportID = detail.ReportID;

DELETE report
FROM dbo.GoreTex_Waterproof_Report report
JOIN @InvalidReports invalid ON invalid.ReportID = report.ReportID;

ALTER TABLE dbo.GoreTex_Waterproof_Report ALTER COLUMN DepNo nvarchar(50) NOT NULL;
ALTER TABLE dbo.GoreTex_Waterproof_Report ALTER COLUMN StyleName nvarchar(100) NOT NULL;
ALTER TABLE dbo.GoreTex_Waterproof_Report_Revision ALTER COLUMN DepNo nvarchar(50) NULL;

/* Restore every legacy line value exactly, including full-width characters. */
UPDATE report
SET report.DepNo = legacy.line
FROM dbo.GoreTex_Waterproof_Report report
JOIN dbo.gore_tex_waterproof_forms legacy
  ON legacy.style_name = report.StyleName
 AND legacy.inspection_date = report.InspectionDate;

UPDATE revision
SET revision.DepNo = report.DepNo
FROM dbo.GoreTex_Waterproof_Report_Revision revision
JOIN dbo.GoreTex_Waterproof_Report report ON report.ReportID = revision.ReportID;

IF NOT EXISTS (
    SELECT 1 FROM sys.indexes
    WHERE object_id = OBJECT_ID('dbo.GoreTex_Waterproof_Report')
      AND name = 'UQ_GoreTex_Waterproof_Report_DepNo_StyleName'
)
    CREATE UNIQUE INDEX UQ_GoreTex_Waterproof_Report_DepNo_StyleName
        ON dbo.GoreTex_Waterproof_Report (DepNo, StyleName);

COMMIT TRANSACTION;
