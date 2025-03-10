using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace pasty
{
    public static class Constants
    {
		public const string Url = "https://127.0.0.1:8000";
		public static HttpClient Conn;
		public const string DatabaseFilename = "Credentials.db3";

		public const SQLite.SQLiteOpenFlags Flags =
			// open the database in read/write mode
			SQLite.SQLiteOpenFlags.ReadWrite |
			// create the database if it doesn't exist
			SQLite.SQLiteOpenFlags.Create |
			// enable multi-threaded database access
			SQLite.SQLiteOpenFlags.SharedCache |
			SQLite.SQLiteOpenFlags.ProtectionCompleteUntilFirstUserAuthentication;

		public static string DatabasePath =>
			Path.Combine(FileSystem.AppDataDirectory, DatabaseFilename);
	}
}
