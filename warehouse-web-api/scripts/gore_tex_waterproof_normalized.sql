SET XACT_ABORT ON;
GO

BEGIN TRANSACTION;

/* Normalize names while the new tables are still empty. */
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'DEPNO') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.DEPNO', 'DepNo', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'DEPNAME') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.DEPNAME', 'DepName', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'XIAMING') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.XIAMING', 'StyleName', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Inspection_date') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Inspection_date', 'InspectionDate', 'COLUMN';

/* Line codes can contain Vietnamese and full-width characters. */
ALTER TABLE dbo.GoreTex_Waterproof_Report ALTER COLUMN DepNo nvarchar(50) NOT NULL;
ALTER TABLE dbo.GoreTex_Waterproof_Report ALTER COLUMN StyleName nvarchar(100) NOT NULL;

IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Total_inspecttion_time_0830') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Total_inspecttion_time_0830', 'TotalInspection0830', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Total_inspecttion_time_0930') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Total_inspecttion_time_0930', 'TotalInspection0930', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Total_inspecttion_time_1030') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Total_inspecttion_time_1030', 'TotalInspection1030', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Total_inspecttion_time_1130') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Total_inspecttion_time_1130', 'TotalInspection1130', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Total_inspecttion_time_1330') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Total_inspecttion_time_1330', 'TotalInspection1330', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Total_inspecttion_time_1430') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Total_inspecttion_time_1430', 'TotalInspection1430', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Total_inspecttion_time_1530') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Total_inspecttion_time_1530', 'TotalInspection1530', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Total_inspecttion_time_1630') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Total_inspecttion_time_1630', 'TotalInspection1630', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'Total_inspecttion_overtime') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.Total_inspecttion_overtime', 'TotalInspectionOvertime', 'COLUMN';

IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'lapbieu') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.lapbieu', 'PreparedBy', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'canbohientruong') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.canbohientruong', 'FieldOfficer', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'canbokcs') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.canbokcs', 'KCSOfficer', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'canbosx') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.canbosx', 'ProductionSupervisor', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'USERID') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.USERID', 'UpdatedBy', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'USERDATE') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report.USERDATE', 'UpdatedAt', 'COLUMN';

IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Time_0830') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Time_0830', 'Time0830', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Time_0930') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Time_0930', 'Time0930', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Time_1030') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Time_1030', 'Time1030', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Time_1130') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Time_1130', 'Time1130', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Time_1330') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Time_1330', 'Time1330', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Time_1430') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Time_1430', 'Time1430', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Time_1530') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Time_1530', 'Time1530', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Time_1630') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Time_1630', 'Time1630', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Total_defect') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Total_defect', 'TotalDefect', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Defect_rate') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Defect_rate', 'DefectRate', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'Shoe_side') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.Shoe_side', 'ShoeSide', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'USERID') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.USERID', 'UpdatedBy', 'COLUMN';
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'USERDATE') IS NOT NULL
    EXEC sys.sp_rename 'dbo.GoreTex_Waterproof_Report_Detail.USERDATE', 'UpdatedAt', 'COLUMN';

/* Audit columns and optimistic-concurrency tokens. */
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'CreatedBy') IS NULL
    ALTER TABLE dbo.GoreTex_Waterproof_Report ADD CreatedBy varchar(50) NULL;
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'CreatedAt') IS NULL
    ALTER TABLE dbo.GoreTex_Waterproof_Report ADD CreatedAt datetimeoffset(0) NOT NULL
        CONSTRAINT DF_GoreTex_Waterproof_Report_CreatedAt DEFAULT SYSDATETIMEOFFSET();
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report', 'RowVersion') IS NULL
    ALTER TABLE dbo.GoreTex_Waterproof_Report ADD RowVersion rowversion NOT NULL;

IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'CreatedBy') IS NULL
    ALTER TABLE dbo.GoreTex_Waterproof_Report_Detail ADD CreatedBy varchar(50) NULL;
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'CreatedAt') IS NULL
    ALTER TABLE dbo.GoreTex_Waterproof_Report_Detail ADD CreatedAt datetimeoffset(0) NOT NULL
        CONSTRAINT DF_GoreTex_Waterproof_Report_Detail_CreatedAt DEFAULT SYSDATETIMEOFFSET();
IF COL_LENGTH('dbo.GoreTex_Waterproof_Report_Detail', 'RowVersion') IS NULL
    ALTER TABLE dbo.GoreTex_Waterproof_Report_Detail ADD RowVersion rowversion NOT NULL;

ALTER TABLE dbo.GoreTex_Waterproof_Report ALTER COLUMN UpdatedAt datetimeoffset(0) NULL;
ALTER TABLE dbo.GoreTex_Waterproof_Report_Detail ALTER COLUMN UpdatedAt datetimeoffset(0) NULL;

IF NOT EXISTS (SELECT 1 FROM sys.default_constraints WHERE name = 'DF_GoreTex_Waterproof_Report_UpdatedAt')
    ALTER TABLE dbo.GoreTex_Waterproof_Report ADD CONSTRAINT DF_GoreTex_Waterproof_Report_UpdatedAt DEFAULT SYSDATETIMEOFFSET() FOR UpdatedAt;
IF NOT EXISTS (SELECT 1 FROM sys.default_constraints WHERE name = 'DF_GoreTex_Waterproof_Report_Detail_UpdatedAt')
    ALTER TABLE dbo.GoreTex_Waterproof_Report_Detail ADD CONSTRAINT DF_GoreTex_Waterproof_Report_Detail_UpdatedAt DEFAULT SYSDATETIMEOFFSET() FOR UpdatedAt;

/* SQL Server does not create an index automatically for a foreign key. */
IF NOT EXISTS (
    SELECT 1 FROM sys.indexes
    WHERE object_id = OBJECT_ID('dbo.GoreTex_Waterproof_Report_Detail')
      AND name = 'IX_GoreTex_Waterproof_Report_Detail_ReportID'
)
    CREATE INDEX IX_GoreTex_Waterproof_Report_Detail_ReportID
        ON dbo.GoreTex_Waterproof_Report_Detail (ReportID);

/* Full normalized snapshots used by the web revision viewer. */
IF OBJECT_ID('dbo.GoreTex_Waterproof_Report_Revision', 'U') IS NULL
BEGIN
    CREATE TABLE dbo.GoreTex_Waterproof_Report_Revision (
        RevisionID bigint IDENTITY(1,1) NOT NULL
            CONSTRAINT PK_GoreTex_Waterproof_Report_Revision PRIMARY KEY,
        ReportID varchar(50) NOT NULL,
        VersionNo int NOT NULL,
        DepNo nvarchar(50) NULL,
        DepName nvarchar(100) NULL,
        StyleName nvarchar(100) NULL,
        InspectionDate date NULL,
        TotalInspection0830 int NULL,
        TotalInspection0930 int NULL,
        TotalInspection1030 int NULL,
        TotalInspection1130 int NULL,
        TotalInspection1330 int NULL,
        TotalInspection1430 int NULL,
        TotalInspection1530 int NULL,
        TotalInspection1630 int NULL,
        TotalInspectionOvertime int NULL,
        PreparedBy nvarchar(50) NULL,
        FieldOfficer nvarchar(50) NULL,
        KCSOfficer nvarchar(50) NULL,
        ProductionSupervisor nvarchar(50) NULL,
        CreatedBy varchar(50) NULL,
        CreatedAt datetimeoffset(0) NULL,
        UpdatedBy varchar(50) NULL,
        UpdatedAt datetimeoffset(0) NULL,
        ArchivedBy varchar(50) NULL,
        ArchivedAt datetimeoffset(0) NOT NULL
            CONSTRAINT DF_GoreTex_Waterproof_Report_Revision_ArchivedAt DEFAULT SYSDATETIMEOFFSET(),
        CONSTRAINT UQ_GoreTex_Waterproof_Report_Revision UNIQUE (ReportID, VersionNo),
        CONSTRAINT FK_GoreTex_Waterproof_Report_Revision_ReportID FOREIGN KEY (ReportID)
            REFERENCES dbo.GoreTex_Waterproof_Report (ReportID)
    );
END;

IF OBJECT_ID('dbo.GoreTex_Waterproof_Report_Detail_Revision', 'U') IS NULL
BEGIN
    CREATE TABLE dbo.GoreTex_Waterproof_Report_Detail_Revision (
        RevisionDetailID bigint IDENTITY(1,1) NOT NULL
            CONSTRAINT PK_GoreTex_Waterproof_Report_Detail_Revision PRIMARY KEY,
        RevisionID bigint NOT NULL,
        DetailID int NULL,
        ReportID varchar(50) NOT NULL,
        Defect nvarchar(100) NULL,
        Time0830 int NULL,
        Time0930 int NULL,
        Time1030 int NULL,
        Time1130 int NULL,
        Time1330 int NULL,
        Time1430 int NULL,
        Time1530 int NULL,
        Time1630 int NULL,
        Overtime int NULL,
        TotalDefect int NULL,
        DefectRate float NULL,
        ShoeSide tinyint NULL,
        CreatedBy varchar(50) NULL,
        CreatedAt datetimeoffset(0) NULL,
        UpdatedBy varchar(50) NULL,
        UpdatedAt datetimeoffset(0) NULL,
        CONSTRAINT FK_GoreTex_Waterproof_Report_Detail_Revision_RevisionID FOREIGN KEY (RevisionID)
            REFERENCES dbo.GoreTex_Waterproof_Report_Revision (RevisionID) ON DELETE CASCADE
    );

    CREATE INDEX IX_GoreTex_Waterproof_Report_Detail_Revision_RevisionID
        ON dbo.GoreTex_Waterproof_Report_Detail_Revision (RevisionID);
END;

IF NOT EXISTS (
    SELECT 1 FROM sys.indexes
    WHERE object_id = OBJECT_ID('dbo.GoreTex_Waterproof_Report')
      AND name = 'UQ_GoreTex_Waterproof_Report_DepNo_StyleName'
)
    CREATE UNIQUE INDEX UQ_GoreTex_Waterproof_Report_DepNo_StyleName
        ON dbo.GoreTex_Waterproof_Report (DepNo, StyleName);

COMMIT TRANSACTION;
GO
