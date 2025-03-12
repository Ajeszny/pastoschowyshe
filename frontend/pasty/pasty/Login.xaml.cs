using pasty.Models;
using System.Text.Json;
using static System.Net.Mime.MediaTypeNames;

namespace pasty;

public partial class LoginPage : ContentPage
{
	public LoginPage()
	{
		InitializeComponent();
	}

	private async void Button_Pressed(object sender, EventArgs e)
	{
		var creds = Credentials.Text;
		var pword = Password.Text;
		try
		{
			var content = new StringContent($"{{\"Credentials\": \"{creds}\", \"Password\": \"{pword}\"}}");
			var response = await Constants.Conn.PostAsync(Constants.Url + "/login", content);
			var t = JsonSerializer.Deserialize<Auth_token>(response.Content.ReadAsStream());
			Constants.db.save_credentials(t);
			await Navigation.PopAsync(true);//return to the main page
		}
		catch (HttpRequestException exception)
		{
			await DisplayAlert("", "Unable to connect to backend: " + exception.Message, "OK");
		}
	}
}