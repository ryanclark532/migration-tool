/*
  Add AccountName from the CdrRequestOpenBankingStatementHistory table
*/
IF NOT EXISTS (
  SELECT *
  FROM sys.columns
  WHERE Name = N'AccountName'
  AND Object_ID = Object_ID(N'[dbo].[CdrRequestOpenBankingStatementHistory]')
)
BEGIN
  ALTER TABLE [dbo].[CdrRequestOpenBankingStatementHistory]
  ADD [AccountName] VARCHAR(256) NULL;
  ADD [AccountNumber] VARCHAR(256) NULL;
END

