#### 1. 为什么需要抽象类?
answer: 抽出一个抽象父类, 为了让无数个类中的**共有逻辑复用**, 不用在每个类中再复制一份代码, 直接继承这个抽象父类就可以      
* 为了代码复用的话, 直接用普通类不就行了, 为什么要用抽象父类?       
answer: 因为所有这些有公共逻辑的类, 往往都是一类业务的实现, 他们不仅有共有逻辑, 还有为了同一个业务目标的**个性化实现逻辑**. 这个个性化逻辑就是在新扩展子类时需要复写的方法. 通过在父类中定义抽象方法, 即能在扩展子类时强制实现(避免扩展子类时忘记腹泻方法), 又能在调用时满足多态

#### 2. 为什么需要接口?
answer: 因为接口是完全没有实现逻辑的, 所以“接口”就是一组“协议”或者“约定”, 相当于一个类的模板, 是对一段逻辑的完全抽象. 子类通过"填空"的方式去实现一段具体逻辑. 同时也增强了代码可读性, 先读接口提出主干脉络

#### 3. 抽象类和接口的设计初衷  
从类的继承层次上来看，抽象类是一种自下而上的设计思路，先有子类的代码重复，然后再抽象成上层的父类（也就是抽象类）。  
    而接口正好相反，它是一种自上而下的设计思路。我们在编程的时候，一般都是先设计接口，再去考虑具体的实现

#### 4. 如何做到上游系统面向接口编程   
* "面向接口编程" 的设计初衷是，将接口和实现相分离，封装不稳定的实现，暴露稳定的接口。上游系统面向接口而非实现编程，不依赖不稳定的实现细节，这样当实现发生变化的时候，上游系统的代码基本上不需要做改动，以此来降低代码间的耦合性，提高代码的扩展性。    
* 不需要抽象出接口的场景: 从这个设计初衷上来看，如果在我们的业务场景中，某个功能只有一种实现方式，未来也不可能被其他实现方式替换，那我们就没有必要为其设计接口，也没有必要基于接口编程，直接使用实现类就可以了。除此之外，越是不稳定的系统，我们越是要在代码的扩展性、维护性上下功夫。相反，如果某个系统特别稳定，在开发完之后，基本上不需要做维护，那我们就没有必要为其扩展性，投入不必要的开发时间。

#### 5. 继承还是组合? 
假设我们要设计一个关于鸟的类。我们将“鸟类”定义为一个抽象类 AbstractBird。所有更细分的鸟，比如麻雀、鸽子、乌鸦等，都继承这个抽象类。  
那 AbstractBird 中要写入什么共有方法呢? 写一个关于飞行的函数 fly()? 不行, 因为有的鸟不会飞, 比如鸵鸟, 鸵鸟类里就不能有 fly() 这个方法, 所以 fly() 不能定义在父类中. 
同理, 关于鸣叫的方法 tweet(), 关于下蛋的方法 egg() 都不能写在父类中. 因为有的鸟不会下蛋, 不会鸣叫.造成这种问题的根本原因是, 对于鸟这类对象, 具体的鸟没有完全相同的一段逻辑, 而是有相同的维度来区分. 因此, 若用继承的方法实现 "飞行,下蛋,鸣叫" 这三种维度, 要三层子类, 共15个抽象类:  
```                         
                                                AbstractBird
                                                     ||
                                AbstractFlyableBird     AbstractUnFlyableBird
                                        ||                       ||
                     可飞行可鸣叫的鸟 , 可飞行不可鸣叫的鸟     不可飞行可鸣叫的鸟 , 不可飞行不可鸣叫的鸟   
                           ||            ||                      ||             ||
      可飞行可鸣叫可下蛋的鸟, 可飞行可鸣叫不可下蛋的鸟      ... ...                      ... (第三层8个抽象子类)   
```
这就体现了继承实现维度, 其实是一种排列组合的类定义方式, 造成父类爆炸: 一方面可读性降低: 要高清一个具体实现类的逻辑, 必须阅读父类的代码、父类的父类的代码 …… 一直追溯到最顶层父类的代码。(2) 另一方面，这也破坏了类的封装特性，将父类的实现细节暴露给了子类。子类的实现依赖父类的实现，两者高度耦合，一旦父类代码修改，就会影响所有子类的逻辑。

1. 现在, 如果我们用接口代替抽象父类, 把三个维度 "飞行,下蛋,鸣叫" 定义成三个接口, 具体实现类去实现这3个接口
    ```java
    public interface Flyable {
      void fly();
    }
    public interface Tweetable {
      void tweet();
    }
    public interface EggLayable {
      void layEgg();
    }
    public class Ostrich implements Tweetable, EggLayable {//鸵鸟
      //... 省略其他属性和方法...
      @Override
      public void tweet() { //... }
      @Override
      public void layEgg() { //... }
    }
    public class Sparrow impelents Flyable, Tweetable, EggLayable {//麻雀
      //... 省略其他属性和方法...
      @Override
      public void fly() { //... }
      @Override
      public void tweet() { //... }
      @Override
      public void layEgg() { //... }
    }
    ```
2. 这样就解决了多维度继承带来的父类爆炸, 但也因为接口无法写具体实现, 造成所有会飞的鸟都要复制一遍 fly(), 所有会鸣叫的鸟都复制一遍 tweet(), 所有会下蛋的鸟都复制一遍 egg().所以现在我们针对三个接口再定义三个实现类：实现了 fly() 方法的 FlyAbility 类、实现了 tweet() 方法的 TweetAbility 类、实现了 layEgg() 方法的 EggLayAbility 类。然后，定义具体的 "鸵鸟" 类时(会下蛋会叫,不会飞), 通过实现 Tweetable, EggLayable 接口, 组合 TweetAbility, EggLayAbility 类, 并将鸵鸟的 fly(), tweet() 逻辑委托给 TweetAbility 和 EggLayAbility 类
    ```java
    public interface Flyable {
        void fly()；
    }
    public class FlyAbility implements Flyable {
        @Override
        public void fly() { //... }
        }
    //省略Tweetable/TweetAbility/EggLayable/EggLayAbility
    
        public class Ostrich implements Tweetable, EggLayable {//鸵鸟
            private TweetAbility tweetAbility = new TweetAbility(); //组合
            private EggLayAbility eggLayAbility = new EggLayAbility(); //组合
    
            //... 省略其他属性和方法...
            @Override
            public void tweet() {
                tweetAbility.tweet(); // 委托
            }
    
            @Override
            public void layEgg() {
                eggLayAbility.layEgg(); // 委托
            }
        }
    }
    ```
所以, 编程鼓励多用组合少用继承