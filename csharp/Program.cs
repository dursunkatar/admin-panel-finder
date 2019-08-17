using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Net;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading;
using System.Threading.Tasks;

namespace AdminFinder
{
    class Program
    {
        static volatile object obj = new object();
        static readonly List<string> panels = new List<string>();
        static int threadRunningCount;
        static string url;
        static int index = -1;

        static void Main(string[] args)
        {
            Console.WriteLine();
            Console.WriteLine(" #####################################");
            Console.WriteLine(" #        Admin Panel Finder         #");
            Console.WriteLine(" #-----------------------------------#");
            Console.WriteLine(" #       github.com/dursunkatar      #");
            Console.WriteLine(" #####################################");

            url = args[0];

            var regx = new Regex("http(s)?://([\\w+?\\.\\w+])+([a-zA-Z0-9\\~\\!\\@\\#\\$\\%\\^\\&amp;\\*\\(\\)_\\-\\=\\+\\\\\\/\\?\\.\\:\\;\\'\\,]*)?", RegexOptions.IgnoreCase);
            if (!regx.IsMatch(url))
            {
                Console.WriteLine("\n Invalid url!\n");
                return;
            }

            string err = loadPanels(args[1]);
            if (err != null)
            {
                Console.WriteLine(err);
                return;
            }

            if (!url.EndsWith("/"))
            {
                url = url + "/";
            }

            Console.WriteLine("\n Panel Url Count: " + panels.Count);
            Console.WriteLine("\n Started...");
            int threadCount = threadRunningCount = 10;
            for (int i = 0; i < threadCount; i++)
            {
                index++;
                connect(url, panels[index]);
            }
        }

        static void connect(string _url, string panel)
        {
            var th = new Thread(() =>
            {
                using (var client = new WebClient())
                {
                    client.Encoding = Encoding.UTF8;
                    try
                    {
                        var result = client.DownloadString(_url + panel);
                        if (result.Contains(" type=\"password\" "))
                        {
                            Console.WriteLine("\n Panel found: " + panel);
                            Console.WriteLine("\n Finish");
                            Environment.Exit(0);
                        }
                    }
                    catch { }
                    thisIsNot();
                }
            });
            th.Start();
        }

        static void thisIsNot()
        {
            lock (obj)
            {
                index++;
                threadRunningCount--;
                if (index < panels.Count)
                {
                    threadRunningCount++;
                    connect(url, panels[index]);
                }
                else if (threadRunningCount == 0)
                {
                    Console.WriteLine("\n Finish");
                }
            }
        }

        static string loadPanels(string path)
        {
            if (!File.Exists(path))
            {
                return "\n " + path + " not exists!\n";
            }
            string[] lines = File.ReadAllLines(path, Encoding.Default);

            foreach (var line in lines)
            {
                string _line = line.Trim(' ', '/');
                if (!panels.Contains(_line))
                {
                    panels.Add(_line);
                }
            }
            return null;
        }
    }
}
