using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace pasty.Models
{
    public class Pasta
    {
        public int Id { get; set; }
        public string Name { get; set; }
        public string[] Tags { get; set; }
    }

    public class Pasta_Text:Pasta
    {
        public string Text { get; set; }
    }
}
