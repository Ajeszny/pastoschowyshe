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
		bool menu_active;
		Button[] menu_buttons;

        public MainPage()
        {
			menu_active = false;
			init = false;
			db = Constants.db;//Let's not refactor too much, shall we?
            vm = new MainPageViewModel(new Command<Pasta>(Transition));
            BindingContext = vm;
            InitializeComponent();
			menu_buttons = [Sudomode, Random, Add];
			foreach (var button in menu_buttons)
			{
				button.ScaleTo(0, 0);
			}
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
			//Maybe I'll require user input to quick log in later?
			Constants.token = await Constants.db.get_credentials();
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

		private async void Menu_Pressed(object sender, EventArgs e)
		{
			animate_button();
		}

		private async void animate_button()
		{
			int[] x_pos = [-50, 50, -50];
			int[] y_pos = [-50, -50, 50];
			if (!menu_active)
			{
				for (global::System.Int32 i = 0; i < menu_buttons.Length; i++)
				{
					await menu_buttons[i].TranslateTo(100, 100, 0);
					menu_buttons[i].ScaleTo(1);
					menu_buttons[i].TranslateTo(x_pos[i], y_pos[i]);
				}
				menu_active = true;
			}
			else
			{
				for (global::System.Int32 i = 0; i < menu_buttons.Length; i++)
				{
					menu_buttons[i].ScaleTo(0);
					menu_buttons[i].TranslateTo(50, 50);
					//await menu_buttons[i].TranslateTo(-100, -100, 0);
				}
				menu_active = false;
			}
		}

		private void Sudomode_Pressed(object sender, EventArgs e)
		{
			if (Constants.token is not null)
			{
				db.logout();
				return;
			}
			Navigation.PushAsync(new LoginPage());
		}

		private void Random_Pressed(object sender, EventArgs e)
		{
			var randomiser = new System.Random();
			var pasta = vm.Pasty[randomiser.Next(vm.Pasty.Count)];
			Transition(pasta);
        }
    }

}
