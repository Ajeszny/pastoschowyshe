using pasty.Models;
using pasty.ViewModels;
using System.Text.Json;

namespace pasty
{
    public partial class MainPage : ContentPage
    {
        int count = 0;

        MainPageViewModel vm;
		Database db;

        public MainPage()
        {
			db = new Database();
            vm = new MainPageViewModel(new Command<Pasta>(Transition));
            BindingContext = vm;
            InitializeComponent();
        }

		private async void Transition(Pasta p)
		{
			var text = await db.get_pasta(p.Id);
			if (text is not null)
			{
				await Navigation.PushAsync(new PastaPage(text));
			}
			var message = new HttpRequestMessage(HttpMethod.Get, Constants.Url + "/get_pasta/" + p.Id);
			try
			{
				var response = await Constants.Conn.SendAsync(message);
				text = JsonSerializer.Deserialize<Pasta_Text>(response.Content.ReadAsStream());
				db.add_pasta(text);
				await Navigation.PushAsync(new PastaPage(text));
			}
			catch (HttpRequestException exception)
			{
				await DisplayAlert("", "Unable to connect to backend: " + exception.Message, "OK");
			}
		}

		protected override async void OnAppearing()
		{
			base.OnAppearing();
			var l = await db.get_saved_pastas();
			vm.Pasty = [.. vm.Pasty.Concat(await db.get_saved_pastas())];

			var message = new HttpRequestMessage(HttpMethod.Get, Constants.Url + "/get_pasta_list");
			try
			{
				var response = await Constants.Conn.SendAsync(message);
				var loaded_list = JsonSerializer.Deserialize<List<Pasta>>(response.Content.ReadAsStream());
				var to_add = from pasta in loaded_list where (!(vm.Pasty.Where(p => p.Id == pasta.Id).Any())) select pasta;
				foreach (var pasta in to_add)
				{
					vm.Pasty.Add(pasta);
				}
			} catch (HttpRequestException exception) {
				await DisplayAlert("", "Unable to connect to backend: " + exception.Message, "OK");
			}
            
		}
	}

}
