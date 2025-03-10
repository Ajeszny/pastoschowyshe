using pasty.Models;
using SQLite;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace pasty
{
	class Pasta_DB_Friendly:Pasta_Text//Sqlite sure has crappy ORM ngl
	{
		public string Tag_str { get { if (Tags is null) { return null; } { } return string.Join(" ", Tags); } set { Tags = value.Split(" "); } }
	}
    public class Database
    {
		SQLiteAsyncConnection conn;

		async Task Init()
		{
			if (conn is not null)
			{
				return;
			}
			var path = Path.Combine(FileSystem.AppDataDirectory, "Credentials.db3");
			conn = new SQLiteAsyncConnection(Constants.DatabasePath, Constants.Flags);
			if (conn is null)
			{
				throw new Exception();
			}
			var result = await conn.CreateTableAsync<Pasta_DB_Friendly>();
		}

		public async Task<Pasta_Text> get_pasta(int id)
		{
			await Init();
			var result = await conn.Table<Pasta_DB_Friendly>().Where(p => p.Id == id).FirstOrDefaultAsync();
			return result;
		}

		public async Task<List<Pasta>> get_saved_pastas()
		{
			await Init();
			var result =  await conn.Table<Pasta_DB_Friendly>().ToListAsync();
			return new List<Pasta>(result);
		}

		public async void add_pasta(Pasta_Text p)
		{
			await Init();
			if (await conn.Table<Pasta_DB_Friendly>().Where(pst => pst.Id == p.Id).CountAsync() != 0)
			{
				return;
			}
			var np = new Pasta_DB_Friendly();
			np.Name = p.Name;
			np.Text = p.Text;
			np.Id = p.Id;
			np.Tags = p.Tags;
			await conn.InsertAsync(np);
		}
	}
}
