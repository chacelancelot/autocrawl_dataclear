/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package autocrawl;

import java.io.BufferedReader;
import java.io.File;
import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.util.HashSet;
import java.util.Scanner;
import java.util.Set;

/**
 *
 * @author DELL
 */
//public class NewClass {
//    public static void main(String[] args) throws IOException, InterruptedException {
//        AutoCrawler crawler = new AutoCrawler();
//        String FolderName = "file_url_13";
//
//        Scanner sc = new Scanner(new File("D:\\D_Documents\\Nam_4\\MMHNC\\Code\\AutoCrawl\\src\\data_uuid\\file_url_13.txt"));
//        
//        int i = 0;
//        while(sc.hasNextLine()){
//            crawler.callCrawler(FolderName, sc.nextLine());
//            System.out.println("so dong: " +  i);
//            i++;
//            
//        }                      
//        System.out.println("well done!");
//    }
//}

public class NewClass {
    public static void main(String[] args) throws IOException, InterruptedException {
        AutoCrawler crawler = new AutoCrawler();
        String FolderName = "file_url_22"; 

        InputStream is = new FileInputStream("D:\\D_Documents\\Nam_4\\MMHNC\\Code\\"
                + "AutoCrawl\\src\\data_uuid\\file_url_22.txt");
        InputStreamReader isr = new InputStreamReader(is);
        BufferedReader sc = new BufferedReader(isr);
        
        String k;
        int i = 0;
        while((k = sc.readLine()) != null){
            crawler.callCrawler(FolderName, k);
            i++;
            System.out.println("so dong: " + i);
        }
            
            
        System.out.println("well done!");
        
    }
}
