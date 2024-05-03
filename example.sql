CREATE TABLE [dbo].[CdrRequest] (
    [RequestID] [UNIQUEIDENTIFIER] PRIMARY KEY, -- Unique ID, sent to DOC and used to correlate event callbacks.

    --Reuse migration fields
    [ApplicantFileID] [UNIQUEIDENTIFIER] NULL, --ID of the Applicant File this CDR Request belongs to
    [ConsumerEmail] [VARCHAR](256) NULL, --Consumers email address
    [RepoID] [VARCHAR](256) NULL, --Repository ID of the CDR Request
    [Active] [BIT] NOT NULL, -- If this is the active CDR Request. An trusted advisor and consumer should only have one active CdrRequest
    [TrustedAdvisorEmail] [VARCHAR](256) NULL,

    --TOBE MOVED TO 'CdrRequestArtefacts'
    [ReportPdfUri] [VARCHAR](512) NULL, -- Lfile URL for the finance report in PDF form. TO BE MOVED TO 'CdrRequestArtefacts'
    [ReportXslxUri] [VARCHAR](512) NULL, -- Lfile URL for the finance report in XSLX form. TO BE MOVED TO 'CdrRequestArtefacts'
    [ReportJsonUri] [VARCHAR](512) NULL, -- Lfile URL for the finance report in JSON form. TO BE MOVED TO 'CdrRequestArtefacts'
    [TransactionsJsonUri] [VARCHAR](512) NULL, -- Lfile URL for the transactions JSON document. TO BE MOVED TO 'CdrRequestArtefacts'
    [AccountsMetadataJsonUri] [VARCHAR](512) NULL, --Lfile URL for the Accounts Metadata JSON Document. TO BE MOVED TO 'CdrRequestArtefacts'

    --TO BE MOVED TO 'CdrRequestState'
    [Status] [VARCHAR](32) NOT NULL, -- Current status: pending, inProgress, completed, error, expired, or updating. TO BE MOVED TO 'CdrRequestState'
    [EnterDateTime] [DATETIME2] NOT NULL, -- Datetime when the record was created. TO BE MOVED TO 'CdrRequestState'
    [CdrRequestErrorID] INT NULL, -- ID of the last error logged against the CDR request. Should only be set if the status is `error`. TO BE MOVED TO 'CdrRequestState'
    [ErrorDateTime] [DATETIME2] NULL, -- Datetime of the last error logged against the CDR request. TO BE MOVED TO 'CdrRequestState'
    [CompletedDateTime] [DATETIME2] NULL, -- Datetime when the CDR request status changed to 'completed'. TO BE MOVED TO 'CdrRequestState'
    [ExpirationDateTime] [DATETIME2] NULL, -- Expiration date time for consent capture. TO BE MOVED TO 'CdrRequestState'
    [DisclosureConsentExpiryDateTime] [DATETIME2] NULL -- Expiration of consent where data can be refreshed. TO BE MOVED TO 'CdrRequestState'
  )


  