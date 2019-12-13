import java.io.*;
import java.nio.ByteBuffer;
import java.nio.MappedByteBuffer;
import java.nio.channels.FileChannel;

public class Test1 {

    private static final int FILESIZE = 300 * 1024 * 1024;
    private static final int BUFFERSIZE = 10 * 1024 * 1024;

    /**
     * 1. 逐字节写入30M: 112042ms (112s)
     */
    private void writeByStream(String file) throws IOException {
        FileOutputStream fos = null;
        try {
            fos = new FileOutputStream(file);
//            BufferedOutputStream bos = new BufferedOutputStream(fos,BUFFERSIZE);
            for (int i = 0; i < FILESIZE; i++) {
                fos.write(0);
//                bos.write(0);
            }
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        }finally {
            if (fos != null)
                fos.close();
        }
    }

    /**
     * 2. 使用FileOutputStream.write(byte[])写入, 数组1m, 耗时47ms
     *    该方法同使用BufferedOutputStream写入, BufferedOutputStream底层调用包装的OutputStream.write(byte[])方法
     */
    private void writeByStreamAndByteArray(String file) throws IOException {
        FileOutputStream fos = null;
        try {
            fos = new FileOutputStream(file);
            byte[] buf = new byte[BUFFERSIZE];
            for (int i = 0; i < FILESIZE/BUFFERSIZE; i++) {
                fos.write(buf);
            }
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        }finally {
            if (fos != null)
                fos.close();
        }
    }

    /**
     * 3. 使用NIO的ByteBuffer做缓冲, channel作为写入目的地; 共耗时330ms
     */
    private void writeByByteBuffer(String file) throws IOException {
        FileOutputStream fos = null;
        FileChannel channel = null;
        ByteBuffer buf = ByteBuffer.allocate(BUFFERSIZE);
        buf.put(new byte[BUFFERSIZE]);

        try{
            fos = new FileOutputStream(file);
            channel = fos.getChannel();
            for (int i = 0; i < FILESIZE/BUFFERSIZE; i++) {
                buf.flip();  // flip进入"读模式", limit = position(写入后的position); position=0
                channel.write(buf);
                // buf.clear();  反复读取不要调用clear方法, clear()进入写模式, position=0; limit=capacity
            }
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        } finally {
            if (fos != null)
                fos.close();
            if (channel != null)
                channel.close();
        }
    }

    /**
     * 4. 使用DirectedBuffer直接内存映射写入(allocateDirect) : 耗时221ms
     */
    private void writeByDirectBuffer(String file) throws IOException {
        FileOutputStream fos = null;
        FileChannel channel = null;
        ByteBuffer buf = ByteBuffer.allocateDirect(BUFFERSIZE);
        buf.put(new byte[BUFFERSIZE]);

        try{
            fos = new FileOutputStream(file);
            channel = fos.getChannel();
            for (int i = 0; i < FILESIZE/BUFFERSIZE; i++) {
                buf.flip();  // flip进入"读模式", limit = position(写入后的position); position=0
                channel.write(buf);
                // buf.clear();  反复读取不要调用clear方法, clear()进入写模式, position=0; limit=capacity
            }
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        } finally {
            if (fos != null)
                fos.close();
            if (channel != null)
                channel.close();
        }
    }
    /**
     * 5. 使用MappedByteBuffer直接内存映射写入 : 耗时405ms
     */
    private void writeByMappedByteBuffer(String file) throws IOException {
        RandomAccessFile raFile = null;
        FileChannel channel = null;
        try {
            raFile = new RandomAccessFile(file, "rw");
            channel = raFile.getChannel();
            byte[] bytes = new byte[BUFFERSIZE];
            for (int i = 0; i < FILESIZE/BUFFERSIZE; i++) {
                MappedByteBuffer buffer = channel.map(FileChannel.MapMode.READ_WRITE, 0+i*BUFFERSIZE, BUFFERSIZE);
                buffer.put(bytes);
            }
        } catch (FileNotFoundException e) {
            e.printStackTrace();
        } finally {
            if (raFile!=null)
                raFile.close();
            if (channel==null)
                channel.close();
        }
    }


    public static void main(String[] args) throws IOException {
        long start = System.currentTimeMillis();
        Test1 t = new Test1();

        t.writeByStream("fos.txt");  // 300M已经太久了
//        t.writeByStreamAndByteArray("fos.txt");  // 151ms
//        t.writeByByteBuffer("fos.txt"); // 140ms
//        t.writeByDirectBuffer("fos.txt"); // 124ms
//        t.writeByMappedByteBuffer("fos.txt"); // 243ms

        long end = System.currentTimeMillis();
        System.out.println(end-start + "ms");
    }
}
