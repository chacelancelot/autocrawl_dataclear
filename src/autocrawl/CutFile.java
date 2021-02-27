/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */
package autocrawl;

import java.io.BufferedReader;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.OutputStreamWriter;

/**
 *
 * @author DELL
 */
public class CutFile {
    public static void run1() throws FileNotFoundException, IOException{
        InputStream is = new FileInputStream("D:\\D_Documents\\Nam_4\\MMHNC\\Code\\AutoCrawl\\src\\data_uuid\\file_url_14.txt");
        InputStreamReader isr = new InputStreamReader(is);
        BufferedReader br = new BufferedReader(isr);
        
        OutputStream os = new FileOutputStream("D:\\D_Documents\\Nam_4\\MMHNC\\Code\\AutoCrawl\\src\\data_uuid\\file_url_36.txt");
        OutputStreamWriter osw = new OutputStreamWriter(os);
        int i = 0, j = 4;
        String k;
        while((k = br.readLine()) != null){
            osw.flush();
            if(i == 1000){
                String s = "D:\\D_Documents\\Nam_4\\MMHNC\\Code\\AutoCrawl\\src\\data_uuid\\file_url_"+j+".txt";
                os = new FileOutputStream(s);
                osw = new OutputStreamWriter(os);
                j++;
                i = 0;
            }
            String o = k;
            osw.write(o);
            osw.write("\n");
            i++;
        }
        osw.flush();
        System.out.println("DONE!");
        osw.close();
        br.close();
        isr.close();
        os.close();
    }
    
    public static void main(String[] args) throws FileNotFoundException, IOException {
        run1();
    }
}
