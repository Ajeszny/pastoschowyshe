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
		bool? init;

        public MainPage()
        {
			init = false;
			db = new Database();
            vm = new MainPageViewModel(new Command<Pasta>(Transition));
            BindingContext = vm;
            InitializeComponent();
		}

		private async void Transition(Pasta p)
		{
			this.TranslateTo(-100, 0);
			var text = await db.get_pasta(p.Id);
			if (text is not null)
			{
				await Navigation.PushAsync(new PastaPage(text));
				return;
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
			while (init is null) ;
			if (init == true)
			{
				await this.TranslateTo(0, 0);
				return;
			}
			init = true;
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

		private async void PointerGestureRecognizer_PointerPressed(object sender, PointerEventArgs e)
		{
			var grid = (Grid)sender;
			var child = (Label)grid.Children[0];

			//var a = new Animation(e => { grid.BackgroundColor = Color.FromRgba("a9a9a9"); });
			//a.Commit(this, "huh");
			await child.ScaleTo(0.9);
		}

		private async void PointerGestureRecognizer_PointerReleased(object sender, PointerEventArgs e)
		{
			var grid = (Grid)sender;
			var child = (Label)grid.Children[0];

			//var a = new Animation(e => { grid.BackgroundColor = Color.FromRgba("ffffff"); });
			//a.Commit(this, "huh");
			await child.ScaleTo(1);
		}
	}

}
