using pasty.Models;
using System;
using System.Collections.Generic;
using System.Collections.ObjectModel;
using System.ComponentModel;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Windows.Input;

namespace pasty.ViewModels
{
    public class MainPageViewModel : INotifyPropertyChanged
	{
		ObservableCollection<Pasta> _pasty;
		public ObservableCollection<Pasta> Pasty { get { return _pasty; } set { _pasty = value; OnPropertyChanged(nameof(Pasty)); } }
		public ICommand OnSwiped { get; set; }
		public event PropertyChangedEventHandler? PropertyChanged;
		protected void OnPropertyChanged(string propertyName)
		{
			PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
		}
		public MainPageViewModel(ICommand on_swiped)
		{
			OnSwiped = on_swiped;
			_pasty = new ObservableCollection<Pasta>();
		}
	}
}
