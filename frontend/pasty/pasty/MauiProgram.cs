using Microsoft.Extensions.Logging;

namespace pasty;

public static class MauiProgram
{
	public static MauiApp CreateMauiApp()
	{
		var handler = new HttpClientHandler
		{
			// Always return true to bypass all certificate checks.
			ServerCertificateCustomValidationCallback = (sender, cert, chain, sslPolicyErrors) => true
		};
		Constants.Conn = new HttpClient(handler);
		Constants.db = new Database();
		var builder = MauiApp.CreateBuilder();
		builder
			.UseMauiApp<App>()
			.ConfigureFonts(fonts =>
			{
				fonts.AddFont("OpenSans-Regular.ttf", "OpenSansRegular");
				fonts.AddFont("OpenSans-Semibold.ttf", "OpenSansSemibold");
			});

#if DEBUG
		builder.Logging.AddDebug();
#endif

		return builder.Build();
	}
}
