package util

import java.lang.reflect.Method
import java.util.UUID

import org.apache.avro.generic.GenericData.StringType
import org.apache.spark.sql.SparkSession
import org.apache.spark.sql.catalyst.encoders.ExpressionEncoder
import org.apache.spark.sql.catalyst.expressions.{Expression, ScalaUDF}
import org.apache.spark.sql.catalyst.{FunctionIdentifier, JavaTypeInference, ScalaReflection}
import org.apache.spark.sql.types.DataType

import scala.tools.reflect.ToolBox
import scala.util.Try

object DynamicUDF {
  /**
    * import util._
    * lazy val trie = new MTrie("testWord")
    * def apply(word:String, matchType:Int):Int={
    *   if(matchType == 0) {
    *     trie.equalMatch(word)
    *   } else {
    *     trie.containsMatch(word)
    *   }
    * }
    */
  def register(spark:SparkSession, func:String, name:String) = {
    val (fun, argumentTypes, returnType) = ScalaGenerateFuns(func)
    val inputTypes = argumentTypes.map{
      x=>
        Try{
          ExpressionEncoder.javaBean(x)
        }.toOption
    }.toSeq
    def builder(e: Seq[Expression]) = ScalaUDF(fun, returnType, e, inputTypes, Some(name))

    spark.sessionState.functionRegistry.registerFunction(new FunctionIdentifier(name), builder)
  }

  case class ClassInfo(clazz: Class[_], instance: Any, defaultMethod: Method, methods: Map[String, Method], func:String) {
    def invoke[T](args: Object*): T = {
      defaultMethod.invoke(instance, args: _*).asInstanceOf[T]
    }
  }

  object ClassCreateUtils {
    private val clazzs = new java.util.HashMap[String, ClassInfo]()
    private val classLoader = this.getClass.getClassLoader
    private val toolBox = scala.reflect.runtime.universe.runtimeMirror(classLoader).mkToolBox()
    def apply(func: String): ClassInfo = this.synchronized {
      var clazz = clazzs.get(func)
      if (clazz == null) {
        val (className, classBody) = wrapClass(func)
        val zz = compile(prepareScala(className, classBody))
        val defaultMethod = zz.getDeclaredMethods.filter(_.getName.contains("apply")).head
        val methods = zz.getDeclaredMethods
        clazz = ClassInfo(
          zz,
          zz.newInstance(),
          defaultMethod,
          methods = methods.map { m => (m.getName, m) }.toMap,
          func
        )
        clazzs.put(func, clazz)
      }
      clazz
    }
    def compile(src: String): Class[_] = {
      toolBox.eval(toolBox.parse(src)).asInstanceOf[Class[_]]
    }
    def prepareScala(className: String, classBody: String): String = {
      classBody + "\n" + s"scala.reflect.classTag[$className].runtimeClass"
    }
    def wrapClass(function: String): (String, String) = {
      val className = s"dynamic_class_${UUID.randomUUID().toString.replaceAll("-", "")}"
      val classBody =
        s"""
           |class $className{
           |  $function
           |}
            """.stripMargin
      (className, classBody)
    }
  }

  object ScalaGenerateFuns {

    def apply(func: String) = {
      val (argumentTypes, returnType) = getFunctionReturnType(func)
      val nax = generateFunction(func, argumentTypes.length)
      (nax, argumentTypes, returnType)
    }

    //获取方法的参数类型及返回类型
    private def getFunctionReturnType(func: String) = {
      val classInfo = ClassCreateUtils(func)
      val method = classInfo.defaultMethod
      val dataType =  JavaTypeInference.inferDataType(method.getReturnType)._1
      val paraTypes = method.getParameterTypes
      (paraTypes, dataType)
    }

    //生成22个Function
    def generateFunction(func: String, argumentsNum: Int): AnyRef = {
      lazy val instance = ClassCreateUtils(func).instance
      lazy val method = ClassCreateUtils(func).methods("apply")

      argumentsNum match {
        case 0 => new (() => Any) with Serializable {
          override def apply(): Any = {
            try {
              method.invoke(instance)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 1 => new ((Object) => Any) with Serializable {
          override def apply(v1: Object): Any = {
            try {
              method.invoke(instance, v1)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 2 => new ((Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object): Any = {
            try {
              method.invoke(instance, v1, v2)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }

        case 3 => new ((Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }


        case 4 => new ((Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }

        case 5 => new ((Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }

        case 6 => new ((Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }

        case 7 => new ((Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }

        case 8 => new ((Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }

        case 9 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }

        case 10 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }

        case 11 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }


        case 12 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }

        case 13 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 14 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object, v14: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 15 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object, v14: Object, v15: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 16 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object, v14: Object, v15: Object, v16: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 17 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object, v14: Object, v15: Object, v16: Object, v17: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 18 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object, v14: Object, v15: Object, v16: Object, v17: Object, v18: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 19 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object, v14: Object, v15: Object, v16: Object, v17: Object, v18: Object, v19: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 20 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object, v14: Object, v15: Object, v16: Object, v17: Object, v18: Object, v19: Object, v20: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 21 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object, v14: Object, v15: Object, v16: Object, v17: Object, v18: Object, v19: Object, v20: Object,
                             v21: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
        case 22 => new ((Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object, Object) => Any) with Serializable {
          override def apply(v1: Object, v2: Object, v3: Object, v4: Object, v5: Object, v6: Object, v7: Object, v8: Object, v9: Object, v10: Object,
                             v11: Object, v12: Object, v13: Object, v14: Object, v15: Object, v16: Object, v17: Object, v18: Object, v19: Object, v20: Object,
                             v21: Object, v22: Object): Any = {
            try {
              method.invoke(instance, v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, v11, v12, v13, v14, v15, v16, v17, v18, v19, v20, v21, v22)
            } catch {
              case e: Exception =>
                e.printStackTrace()
                null
            }
          }
        }
      }
    }
  }
}
