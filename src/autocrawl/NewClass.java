/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package autocrawl;

import java.io.File;
import java.io.IOException;
import java.util.HashSet;
import java.util.Scanner;
import java.util.Set;

/**
 *
 * @author DELL
 */
public class NewClass {
    public static void main(String[] args) throws IOException, InterruptedException {
        AutoCrawler crawler = new AutoCrawler();
        String FolderName = "file_url1";

        Scanner sc = new Scanner(new File("D:\\D_Documents\\Nam_4\\MMHNC\\Code\\AutoCrawl\\src\\data_uuid\\file_url.txt"));
        
        while(sc.hasNextLine())
            crawler.callCrawler(FolderName, sc.nextLine());
            
        
    }
}
