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
        static string[] panels;
        static int threadRunningCount;
        static string url;
        static int index = -1;

        static void Main(string[] args)
        {
            Console.WriteLine();
            Console.WriteLine(
   @"    ##################################### 
    #        Admin Panel Finder         #
    #-----------------------------------#
    #       Author: Dursun Katar        #
    #-----------------------------------#
    #       github.com/dursunkatar      #
    #####################################");

            url = args[0];

            var regx = new Regex("http(s)?://([\\w+?\\.\\w+])+([a-zA-Z0-9\\~\\!\\@\\#\\$\\%\\^\\&amp;\\*\\(\\)_\\-\\=\\+\\\\\\/\\?\\.\\:\\;\\'\\,]*)?", RegexOptions.IgnoreCase);
            if (!regx.IsMatch(url))
            {
                Console.WriteLine("\n Invalid url!\n");
                return;
            }

            if (!File.Exists(args[1]))
            {
                Console.WriteLine("\n File not exists!\n");
                return;
            }

            if (!url.EndsWith("/"))
            {
                url = url + "/";
            }

            loadPanels(args[1]);

            Console.WriteLine("\n Panel Url Count: " + panels.Length);
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
                            index = panels.Length;
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
                if (index < panels.Length)
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

        static void loadPanels(string path)
        {
            string[] lines = File.ReadAllLines(path, Encoding.Default);
            var list = new List<string>();
            foreach (var line in lines)
            {
                string _line = line.Trim(' ', '/');
                if (!list.Contains(_line))
                {
                    list.Add(_line);
                }
            }
            panels = list.ToArray();
        }
    }
}
