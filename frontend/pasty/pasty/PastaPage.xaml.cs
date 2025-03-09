using pasty.Models;
using pasty.ViewModels;

namespace pasty;

public partial class PastaPage : ContentPage
{
	PastaViewModel vm;
	public PastaPage(Pasta_Text pasta)
	{
		vm = new PastaViewModel(pasta);
		BindingContext = vm;
		InitializeComponent();
	}
}