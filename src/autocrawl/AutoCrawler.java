/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package autocrawl;


import java.io.BufferedReader;
import java.io.File;
import java.io.IOException;
import java.io.InputStreamReader;
/**
 *
 * @author DELL
 */
public class AutoCrawler {
    private ProcessBuilder builder;

    public AutoCrawler() throws IOException, InterruptedException {
        builder = new ProcessBuilder();
        builder.directory(new File("D:\\D_Documents\\Nam_4\\MMHNC\\Data_clean\\app.any-crawler-master"));
        builder.command("powershell.exe", "/c", "go build -o app.any-crawler.exe");
        Process process = builder.start();
//        process.waitFor();
    }
    
    
    
    

    public int callCrawler(String  fileName, String id) throws IOException, InterruptedException {
//        ProcessBuilder builder = new ProcessBuilder();
//        builder.directory(new File("D:\\D_Documents\\Nam_4\\MMHNC\\Data_clean\\app.any-crawler-master"));
//        builder.command("powershell.exe", "/c", "go build -o app.any-crawler.exe");
//        builder.directory(new File("D:\\D_Documents\\Nam_4\\MMHNC\\Data_clean\\app.any-crawler-master"));
        builder.command("powershell.exe", "/c", "./app.any-crawler.exe -i "+ fileName + " " + id);
        Process process = builder.start();
        BufferedReader reader
                = new BufferedReader(new InputStreamReader(process.getInputStream()));

        String line;
        while ((line = reader.readLine()) != null) {
            System.out.println(line);
        }

//        int exitCode = process.waitFor();
//        System.out.println("\nExited with error code : " + exitCode);
        return process.waitFor();
     }
   
    
}

/*import os

# os.system('cmd /c "cd ../Data_clean/app.any-crawler-master"' )

os.system('cmd /c "go build -o app.any-crawler.exe"')
#
# fs = open(r"abc.txt", 'r')
#
# listID = fs.readlines();
#
# for id in listID:
#     cmd = "./app.any-crawler.exe -i filedll " + id
#     os.system('cmd /k' + cmd)

*/


