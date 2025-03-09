using pasty.Models;
using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace pasty.ViewModels
{
    public class PastaViewModel : INotifyPropertyChanged
	{
		Pasta_Text _body { get; set; }

		public Pasta_Text Body { get { return _body; } set { _body = value; } }

		public PastaViewModel(Pasta_Text body)
		{
			Body = body;
		}

		public event PropertyChangedEventHandler? PropertyChanged;
		protected void OnPropertyChanged(string propertyName)
		{
			PropertyChanged?.Invoke(this, new PropertyChangedEventArgs(propertyName));
		}
	}
}
