
IF NOT EXISTS (SELECT * FROM dbo.sysobjects   WHERE id = OBJECT_ID(N'Employees') AND OBJECTPROPERTY(id, N'IsTable') = 1 )
	BEGIN
		CREATE TABLE Employees (
			ID INT IDENTITY(1,1) PRIMARY KEY,
			FirstName VARCHAR(256),
			Lastname VARCHAR(256)
		);
	END
